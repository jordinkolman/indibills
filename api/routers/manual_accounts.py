from fastapi import APIRouter, Depends
from fastapi.responses import JSONResponse
from sqlalchemy.exc import NoResultFound
from sqlmodel import Session, select

from ..dependencies import get_session
from ..models import ManualAccountCreate, ManualAccount, ManualAccountUpdate, Message, User


router = APIRouter(
    prefix="/manual_accounts",
    tags=["manual_accounts"],
    responses={404: {"model": Message}}
)


@router.get("/user/{user_id}", response_model=list[ManualAccount])
async def get_manual_accounts_by_user_id(*, session: Session = Depends(get_session), user_id: int):
    user = session.get(User, user_id)
    if not user:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    return session.exec(select(ManualAccount).where(ManualAccount.user_id == user.id)).all()

@router.get("/account/{account_id}", response_model=ManualAccount)
async def get_manual_account_by_id(*, session: Session = Depends(get_session), account_id: int):
    try:
        account = session.exec(select(ManualAccount).where(ManualAccount.id == account_id)).one()
    except NoResultFound:
        return JSONResponse(status_code=404, content={"message": "Account Not Found"})
    return account

@router.post("/", response_model=ManualAccount)
async def create_manual_account(*, session: Session = Depends(get_session), account: ManualAccountCreate):
    db_account = ManualAccount.model_validate(account)
    session.add(db_account)
    session.commit()
    session.refresh(db_account)
    return db_account

@router.patch("/account/{account_id}", response_model=ManualAccount)
async def update_manual_account(*, session: Session = Depends(get_session), account_id: int, account: ManualAccountUpdate):
    db_account = session.get(ManualAccount, account_id)
    if not db_account:
        return JSONResponse(status_code=404, content={"message": "Account Not Found"})
    account_data = account.model_dump(exclude_unset=True)
    db_account.sqlmodel_update(account_data)
    session.add(db_account)
    session.commit()
    session.refresh(db_account)
    return db_account

@router.delete("/account/{account_id}")
async def delete_manual_account(*, session: Session = Depends(get_session), account_id: int):
    account = session.get(ManualAccount, account_id)
    if not account:
        return JSONResponse(status_code=404, content={"message": "Account Not Found"})
    session.delete(account)
    session.commit()
    return{"ok": True}
