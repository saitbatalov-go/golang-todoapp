


func (s *UsersService) PatchUser(ctx context.Context, id int, user domain.UserPatch) (domain.User, error) {

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(user); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	user, err = s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}
	return user, nil
}