import uuid
from typing import Sequence

import pytest
from returns.maybe import Nothing
from returns.result import Failure, Success
from sqlalchemy.ext.asyncio import AsyncSession

from annotatex.models.organization import Organization, OrganizationKind
from annotatex.repositories.organization_repository import (
    OrganizationRepository,
)


@pytest.fixture()
def org_repo(session: AsyncSession) -> OrganizationRepository:
    return OrganizationRepository(session=session)


async def _create(
    repo: OrganizationRepository,
    name: str = "acme",
    kind: OrganizationKind = OrganizationKind.client,
    is_active: bool = True,
) -> Organization:
    res = await repo.create_org(name=name, kind=kind, is_active=is_active)
    assert isinstance(res, Success)
    return res.unwrap()


class TestOrganizationRepository:
    async def test_create_org_success(self, org_repo: OrganizationRepository):
        org = await _create(org_repo, name="alpha", kind=OrganizationKind.vendor)
        assert isinstance(org.id, uuid.UUID)
        assert org.name == "alpha"
        assert org.kind == OrganizationKind.vendor
        assert org.is_active is True

    async def test_create_org_duplicate_name_failure(
        self, org_repo: OrganizationRepository
    ):
        await _create(org_repo, name="dup")
        res2 = await org_repo.create_org(name="dup", kind=OrganizationKind.client)

        assert isinstance(res2, Failure)
        assert "IntegrityError" in type(res2.failure()).__name__

    async def test_get_by_id_and_name(self, org_repo: OrganizationRepository):
        org = await _create(org_repo, name="bykey", kind=OrganizationKind.internal)
        maybe_by_id = await org_repo.get_by_id(org.id)
        assert maybe_by_id != Nothing
        assert maybe_by_id.unwrap().name == "bykey"

        maybe_by_name = await org_repo.get_by_name("bykey")
        assert maybe_by_name != Nothing
        assert maybe_by_name.unwrap().id == org.id

        # miss
        maybe_missing = await org_repo.get_by_id(uuid.uuid4())
        assert maybe_missing == Nothing

    async def test_list_active_and_all(self, org_repo: OrganizationRepository):
        await _create(org_repo, name="a1", is_active=True)
        await _create(org_repo, name="a2", is_active=False)
        await _create(org_repo, name="a3", is_active=True)

        active: Sequence[Organization] = await org_repo.list_active()
        all_items: Sequence[Organization] = await org_repo.list_all()

        assert {o.name for o in active} == {"a1", "a3"}
        assert {o.name for o in all_items} >= {"a1", "a2", "a3"}

    async def test_update_org_partial_fields(self, org_repo: OrganizationRepository):
        org = await _create(
            org_repo, name="u1", kind=OrganizationKind.client, is_active=True
        )

        res = await org_repo.update_org(
            org.id,
            name="u1-renamed",
            kind=OrganizationKind.vendor,
        )
        assert isinstance(res, Success)
        updated = res.unwrap()
        assert updated.id == org.id
        assert updated.name == "u1-renamed"
        assert updated.kind == OrganizationKind.vendor
        assert updated.is_active is True  # 未修改的字段保持不变

    async def test_update_org_noop_returns_current(
        self, org_repo: OrganizationRepository
    ):
        org = await _create(org_repo, name="noop")

        res = await org_repo.update_org(org.id)
        assert isinstance(res, Success)

        same = res.unwrap()
        assert same.id == org.id
        assert same.name == "noop"

    async def test_update_org_noop_not_found_returns_failure(
        self, org_repo: OrganizationRepository
    ):
        res = await org_repo.update_org(uuid.uuid4())
        assert isinstance(res, Failure)
        assert "Not found" in str(res.failure())

    async def test_set_active_activate_deactivate(
        self, org_repo: OrganizationRepository
    ):
        org = await _create(org_repo, name="toggle", is_active=True)

        res1 = await org_repo.deactivate(org.id)
        assert isinstance(res1, Success)
        assert res1.unwrap().is_active is False

        res2 = await org_repo.activate(org.id)
        assert isinstance(res2, Success)
        assert res2.unwrap().is_active is True

        res3 = await org_repo.set_active(org.id, active=False)
        assert isinstance(res3, Success)
        assert res3.unwrap().is_active is False

    async def test_soft_delete(self, org_repo: OrganizationRepository):
        org = await _create(org_repo, name="softy", is_active=True)
        res_true = await org_repo.soft_delete(org.id)
        assert isinstance(res_true, Success)
        assert res_true.unwrap() is True

        res_false = await org_repo.soft_delete(uuid.uuid4())
        assert isinstance(res_false, Failure)

    async def test_bulk_deactivate(self, org_repo: OrganizationRepository):
        o1 = await _create(org_repo, name="b1", is_active=True)
        o2 = await _create(org_repo, name="b2", is_active=True)
        _ = await _create(org_repo, name="b3", is_active=True)

        res = await org_repo.bulk_deactivate([o1.id, o2.id])
        assert isinstance(res, Success)
        assert res.unwrap() == 2

        active = await org_repo.list_active()
        assert {o.name for o in active} == {"b3"}
