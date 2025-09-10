import uuid
from typing import Sequence

from returns.maybe import Maybe
from returns.result import Failure, Result, Success
from sqlalchemy import func, select, update
from sqlalchemy.exc import IntegrityError

from annotatex.models.user import User
from annotatex.repositories.base_repository import BaseRepository


class UserRepository(BaseRepository):
    async def create_user(
        self, *, username: str, is_active: bool = True
    ) -> Result[User, IntegrityError]:
        user = User(username=username, is_active=is_active)
        try:
            async with self.session.begin():
                self.session.add(user)
                await self.session.flush()
            return Success(user)
        except IntegrityError as exc:
            return Failure(exc)

    async def get_user_by_id(self, user_id: uuid.UUID) -> Maybe[User]:
        result = await self.session.execute(select(User).where(User.id == user_id))
        return Maybe.from_optional(result.scalar_one_or_none())

    async def get_user_by_username(self, username: str) -> Maybe[User]:
        result = await self.session.execute(
            select(User).where(User.username == username)
        )
        return Maybe.from_optional(result.scalar_one_or_none())

    async def list_all(self) -> Sequence[User]:
        result = await self.session.execute(
            select(User).order_by(User.created_at.desc())
        )
        return list(result.scalars().all())

    async def list_active(self) -> Sequence[User]:
        result = await self.session.execute(
            select(User)
            .where(User.is_active.is_(True))
            .order_by(User.created_at.desc())
        )
        return list(result.scalars().all())

    async def set_active(
        self, user_id: uuid.UUID, *, is_active: bool
    ) -> Result[User, Exception]:
        try:
            async with self.session.begin():
                stmt = (
                    update(User)
                    .where(User.id == user_id)
                    .values(is_active=is_active, updated_at=func.now())
                    .returning(User)
                )
                result = await self.session.execute(stmt)
                updated = result.scalar_one_or_none()
                if updated is None:
                    raise ValueError("User not found")
            return Success(updated)
        except Exception as exc:
            return Failure(exc)
