from fastapi import APIRouter, Depends, Query
from fastapi.responses import JSONResponse
from sqlalchemy.exc import NoResultFound
from sqlmodel import Session, select

from ..dependencies import get_session
from ..models import Household, HouseholdCreate, HouseholdUpdate, Message

router = APIRouter(
    prefix="/households", tags=["households"], responses={404: {"model": Message}}
)


@router.get("/", response_model=list[Household])
async def get_households(
    *,
    session: Session = Depends(get_session),
    offset: int = 0,
    limit: int = Query(default=100, le=100),
):
    households = session.exec(select(Household).offset(offset).limit(limit)).all()
    if not households:
        return JSONResponse(status_code=404, content={"message": "No Households Found"})
    return households


@router.post("/")
async def create_household(
    *, session: Session = Depends(get_session), household: HouseholdCreate
):
    db_household = Household.model_validate(household)
    session.add(db_household)
    session.commit()
    session.refresh(db_household)
    return db_household


@router.get("/{household_id}", response_model=Household)
async def get_household_by_id(
    *, session: Session = Depends(get_session), household_id: int
):
    try:
        household = session.exec(
            select(Household).where(Household.id == household_id)
        ).one()
    except NoResultFound:
        return JSONResponse(status_code=404, content={"message": "Household Not Found"})
    return household


@router.patch("/{household_id}", response_model=Household)
async def update_household(
    *,
    session: Session = Depends(get_session),
    household_id: int,
    household: HouseholdUpdate,
):
    db_household = session.get(Household, household_id)
    if not db_household:
        return JSONResponse(status_code=404, content={"message": "Household Not Found"})
    household_data = household.model_dump(exclude_unset=True)
    db_household.sqlmodel_update(household_data)
    session.add(db_household)
    session.commit()
    session.refresh(db_household)
    return db_household


@router.delete("/{household_id}")
async def delete_household(
    *, session: Session = Depends(get_session), household_id: int
):
    household = session.get(Household, household_id)
    if not household:
        return JSONResponse(status_code=404, content={"message": "Household Not Found"})
    session.delete(household)
    session.commit()
    return {"ok": True}
