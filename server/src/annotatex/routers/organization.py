from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from returns.result import Failure, Success
from sqlalchemy.ext.asyncio import AsyncSession

from annotatex.db.engine import get_async_session
from annotatex.repositories.organization_repository import OrganizationRepository
from annotatex.schemas.organization_schema import CreateOrganizationSchema

router = APIRouter(prefix="/organization")


@router.post("/")
async def create(
    session: Annotated[AsyncSession, Depends(get_async_session)],
    org: CreateOrganizationSchema,
):
    result_org = await OrganizationRepository(session).create_org(
        name=org.name, kind=org.kind
    )
    match result_org:
        case Success(org_model):
            return org_model
        case Failure(exc):
            raise HTTPException(status_code=400, detail=str(exc))
