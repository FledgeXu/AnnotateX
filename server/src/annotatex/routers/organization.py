from fastapi import APIRouter

router = APIRouter(prefix="/organization")


@router.get("/")
async def get_organizations():
    return "placeholder"
