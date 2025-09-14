import uuid

import pytest
from returns.maybe import Nothing, Some
from returns.result import Failure, Success

from annotatex.repositories.user_repository import UserRepository


@pytest.fixture()
def user_repo(session) -> UserRepository:
    return UserRepository(session)


class TestUserRepository:
    async def test_create_user_success(self, user_repo: UserRepository):
        result = await user_repo.create_user(username="alice")
        assert isinstance(result, Success)
        user = result.unwrap()
        assert isinstance(user.id, uuid.UUID)
        assert user.username == "alice"
        assert user.is_active is True

    async def test_create_user_duplicate(self, user_repo: UserRepository):
        await user_repo.create_user(username="bob")
        result = await user_repo.create_user(username="bob")
        assert isinstance(result, Failure)

    async def test_get_user_by_id_and_username(self, user_repo: UserRepository):
        result = await user_repo.create_user(username="charlie")
        user = result.unwrap()

        by_id = await user_repo.get_user_by_id(user.id)
        assert isinstance(by_id, Some)
        assert by_id.unwrap().username == "charlie"

        by_username = await user_repo.get_user_by_username("charlie")
        assert isinstance(by_username, Some)
        assert by_username.unwrap().id == user.id

        not_found = await user_repo.get_user_by_username("nobody")
        assert not_found == Nothing

    async def test_list_all_and_active(self, user_repo: UserRepository):
        await user_repo.create_user(username="a1", is_active=True)
        await user_repo.create_user(username="a2", is_active=False)
        await user_repo.create_user(username="a3", is_active=True)

        all_users = await user_repo.list_all()
        assert {u.username for u in all_users} == {"a1", "a2", "a3"}

        active_users = await user_repo.list_active()
        assert {u.username for u in active_users} == {"a1", "a3"}

    async def test_set_active_success_and_failure(self, user_repo: UserRepository):
        result = await user_repo.create_user(username="dave")
        user = result.unwrap()

        updated = await user_repo.set_active(user.id, is_active=False)
        assert isinstance(updated, Success)
        assert updated.unwrap().is_active is False

        # not existing id
        fake_id = uuid.uuid4()
        not_found = await user_repo.set_active(fake_id, is_active=True)
        assert isinstance(not_found, Failure)
