from dotenv import load_dotenv
from fastapi import APIRouter, FastAPI

from annotatex.routers.organization import router as organization_router

load_dotenv()
app = FastAPI()


router = APIRouter(prefix="/v1")
router.include_router(organization_router)

app.include_router(router)
