package matchmaking

import (
	"context"
	"fmt"
	"strconv"
)

func (u *Matchmaking) CreateMatchmakingPool(ctx context.Context, userID int) error {
	ctx, s := u.Monitor(ctx, "User: CreateMatchmakingPool")
	defer s.Finish(ctx)

	// TODO: put rank logic here to accomodate interest
	query := `
        SELECT id 
        FROM users 
        WHERE id != $1
          AND id NOT IN (SELECT swiped_id FROM match WHERE swiper_id = $1 AND status = 'active')
          AND id NOT IN (SELECT swiper_id FROM match WHERE swiped_id = $1 AND status = 'active')
    `

	var candidates []int
	err := u.dbR.SelectContext(ctx, &candidates, query, userID)
	if err != nil {
		return fmt.Errorf("[User][CreateMatchmakingPool] failed to fetch eligible users from DB: %w", err)
	}

	swipedUsersKey := fmt.Sprintf(rkSwipedUsers, userID)
	swipedUsers, err := u.redisClient.SMembers(ctx, swipedUsersKey).Result()
	if err != nil {
		return fmt.Errorf("[User][CreateMatchmakingPool] failed to fetch swiped users from Redis: %w", err)
	}

	swipedUserIDs := make(map[int]struct{})
	for _, id := range swipedUsers {
		swipedUserID, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("[User][CreateMatchmakingPool] invalid swiped user ID: %w", err)
		}
		swipedUserIDs[swipedUserID] = struct{}{}
	}

	var finalPool []int
	for _, candidateID := range candidates {
		if _, swiped := swipedUserIDs[candidateID]; !swiped {
			finalPool = append(finalPool, candidateID)
		}
	}

	poolKey := fmt.Sprintf(rkPool, userID)
	err = u.redisClient.SAdd(ctx, poolKey, finalPool).Err()
	if err != nil {
		return fmt.Errorf("[User][CreateMatchmakingPool] failed to add candidates to Redis pool: %w", err)
	}

	err = u.redisClient.Expire(ctx, poolKey, rkPoolExp).Err()
	if err != nil {
		return fmt.Errorf("[User][CreateMatchmakingPool] failed to set expiration on Redis pool: %w", err)
	}

	return nil
}

func (u *Matchmaking) PopMatchmakingPool(ctx context.Context, userID int, count int) ([]int, error) {
	ctx, s := u.Monitor(ctx, "Matchmaking: PopMatchmakingPool")
	defer s.Finish(ctx)

	poolKey := fmt.Sprintf(rkPool, userID)
	poppedUsers, err := u.redisClient.SPopN(ctx, poolKey, int64(count)).Result()
	if err != nil {
		return nil, fmt.Errorf("[Matchmaking][PopMatchmakingPool] failed to pop users from Redis pool: %w", err)
	}

	if len(poppedUsers) == 0 {
		return nil, fmt.Errorf("[Matchmaking][PopMatchmakingPool] no users available in the pool")
	}

	var userIDs []int
	for _, user := range poppedUsers {
		userID, err := strconv.Atoi(user)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in Redis pool: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (u *Matchmaking) InsertMatch(ctx context.Context, swiperID, swipedID int) error {
	ctx, s := u.Monitor(ctx, "Matchmaking: InsertMatch")
	defer s.Finish(ctx)

	swipedUsersKey := fmt.Sprintf(rkSwipedUsers, swiperID)
	swipedUsers, err := u.redisClient.SIsMember(ctx, swipedUsersKey, swipedID).Result()
	if err != nil {
		return fmt.Errorf("[Matchmaking][InsertMatch] failed to check if swiper swiped right: %w", err)
	}

	if !swipedUsers {
		return fmt.Errorf("[Matchmaking][InsertMatch] swiper has not swiped right on the swiped user")
	}

	swipedByUsersKey := fmt.Sprintf(rkSwipedUsers, swipedID)
	swipedByUsers, err := u.redisClient.SIsMember(ctx, swipedByUsersKey, swiperID).Result()
	if err != nil {
		return fmt.Errorf("[Matchmaking][InsertMatch] failed to check if swiped user swiped right on the swiper: %w", err)
	}

	if !swipedByUsers {
		return fmt.Errorf("[Matchmaking][InsertMatch] swiped user has not swiped right on the swiper")
	}

	insertQuery := `
        INSERT INTO match (swiper_id, swiped_id, status) 
        VALUES ($1, $2, 'active')
    `
	_, err = u.dbW.ExecContext(ctx, insertQuery, swiperID, swipedID)
	if err != nil {
		return fmt.Errorf("[Matchmaking][InsertMatch] failed to insert new match: %w", err)
	}

	return nil
}
