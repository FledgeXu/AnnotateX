import uuid
from typing import Iterable, Sequence

from returns.maybe import Maybe
from returns.result import Failure, Result, Success
from sqlalchemy import select, update
from sqlalchemy.exc import IntegrityError, SQLAlchemyError

from annotatex.models.organization import Organization, OrganizationKind
from annotatex.repositories.base_repository import BaseRepository


class OrganizationRepository(BaseRepository):
    async def create_org(
        self, *, name: str, kind: OrganizationKind, is_active: bool = True
    ) -> Result[Organization, IntegrityError]:
        org = Organization(name=name, kind=kind, is_active=is_active)
        try:
            async with self.session.begin():
                self.session.add(org)
                await self.session.flush()
            return Success(org)
        except IntegrityError as exc:
            return Failure(exc)

    async def get_by_id(self, org_id: uuid.UUID) -> Maybe[Organization]:
        async with self.session.begin():
            result = await self.session.execute(
                select(Organization).where(Organization.id == org_id)
            )
            return Maybe.from_optional(result.scalar_one_or_none())

    async def get_by_name(self, name: str) -> Maybe[Organization]:
        async with self.session.begin():
            result = await self.session.execute(
                select(Organization).where(Organization.name == name)
            )
            return Maybe.from_optional(result.scalar_one_or_none())

    async def list_active(self) -> Sequence[Organization]:
        async with self.session.begin():
            result = await self.session.execute(
                select(Organization).where(Organization.is_active.is_(True))
            )
            return list(result.scalars().all())

    async def list_all(self) -> Sequence[Organization]:
        async with self.session.begin():
            result = await self.session.execute(select(Organization))
            return list(result.scalars().all())

    async def update_org(
        self,
        org_id: uuid.UUID,
        *,
        name: str | None = None,
        kind: OrganizationKind | None = None,
        is_active: bool | None = None,
    ) -> Result[Organization, SQLAlchemyError]:
        """
        Partial update. Returns the updated Organization in Success.
        """
        fields: dict = {}
        if name is not None:
            fields["name"] = name
        if kind is not None:
            fields["kind"] = kind
        if is_active is not None:
            fields["is_active"] = is_active

        if not fields:
            # No-op update: fetch current record for convenience
            maybe_org = await self.get_by_id(org_id)
            return maybe_org.map(Success).value_or(
                Failure(SQLAlchemyError("Not found"))
            )

        try:
            async with self.session.begin():
                # Keep updated_at current if your model uses onupdate=func.now()
                # DB will handle it; otherwise you can set it explicitly here.
                stmt = (
                    update(Organization)
                    .where(Organization.id == org_id)
                    .values(**fields)
                    .returning(Organization)
                )
                result = await self.session.execute(stmt)
                org = result.scalar_one_or_none()
                if org is None:
                    raise SQLAlchemyError("Organization not found")
            return Success(org)
        except (IntegrityError, SQLAlchemyError) as exc:
            return Failure(exc)

    async def set_active(
        self, org_id: uuid.UUID, *, active: bool
    ) -> Result[Organization, SQLAlchemyError]:
        return await self.update_org(org_id, is_active=active)

    async def activate(
        self, org_id: uuid.UUID
    ) -> Result[Organization, SQLAlchemyError]:
        return await self.set_active(org_id, active=True)

    async def deactivate(
        self, org_id: uuid.UUID
    ) -> Result[Organization, SQLAlchemyError]:
        return await self.set_active(org_id, active=False)

    async def soft_delete(self, org_id: uuid.UUID) -> Result[bool, SQLAlchemyError]:
        values: dict = {"is_active": False}
        try:
            async with self.session.begin():
                stmt = (
                    update(Organization)
                    .where(Organization.id == org_id)
                    .values(**values)
                    .returning(Organization.id)
                )
                result = await self.session.execute(stmt)
                updated_id = result.scalar_one_or_none()
                if updated_id is None:
                    return Failure(SQLAlchemyError("Upload id not found."))
                return Success(True)
        except SQLAlchemyError as exc:
            return Failure(exc)

    async def bulk_deactivate(
        self, org_ids: Iterable[uuid.UUID]
    ) -> Result[int, SQLAlchemyError]:
        try:
            async with self.session.begin():
                stmt = (
                    update(Organization)
                    .where(Organization.id.in_(list(org_ids)))
                    .values(is_active=False)
                )
                result = await self.session.execute(stmt)
                affected = result.rowcount or 0
            return Success(affected)
        except SQLAlchemyError as exc:
            return Failure(exc)
