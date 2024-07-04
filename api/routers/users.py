from fastapi import APIRouter, Depends, Query
from fastapi.responses import JSONResponse
from passlib.hash import pbkdf2_sha256
from sqlalchemy.exc import NoResultFound
from sqlmodel import Session, select

from ..dependencies import get_session
from ..models import User, UserCreate, UserPublic, UserUpdate, Message

router = APIRouter(
    prefix="/users",
    tags=["users"],
    responses={404: {"model": Message}}
)


@router.post("/", response_model=UserPublic)
async def create_user(*, session: Session = Depends(get_session), user: UserCreate):
    hashed_password = pbkdf2_sha256.hash(user.password)
    extra_data = {"hashed_password": hashed_password}
    db_user = User.model_validate(user, update=extra_data)
    session.add(db_user)
    session.commit()
    session.refresh(db_user)
    return db_user

@router.get("/", response_model=list[UserPublic])
async def get_users(session: Session = Depends(get_session), offset: int = 0, limit: int = Query(default=100, le=100)):
    users = session.exec(select(User).offset(offset).limit(limit)).all()
    if not users:
        return JSONResponse(status_code=404, content={"message": "No Users Found"})
    return users

@router.get("/{username}", response_model=UserPublic)
async def get_user_by_username(*, session: Session = Depends(get_session), username: str):
    try:
        user = session.exec(select(User).where(User.username == username)).one()
    except NoResultFound:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    return user

@router.get("/id/{user_id}", response_model=UserPublic)
async def get_user_by_id(*, session: Session = Depends(get_session), user_id: int):
    try:
        user = session.exec(select(User).where(User.id == user_id)).one()
    except NoResultFound:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    return user

@router.get("/email/{user_email}", response_model=UserPublic)
async def get_user_by_email(*, session: Session = Depends(get_session), user_email: str):
    user = session.exec(select(User).where(User.email == user_email)).one()
    if not user:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    return user

@router.patch("/{user_id}", response_model=UserPublic)
async def update_user(*, session: Session = Depends(get_session), user_id: int, user: UserUpdate):
    db_user = session.get(User, user_id)
    if not db_user:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    user_data = user.model_dump(exclude_unset=True)
    extra_data = {}
    if "password" in user_data:
        password = user_data["password"]
        hashed_password = pbkdf2_sha256.hash(password)
        extra_data["hashed_password"] = hashed_password
    db_user.sqlmodel_update(user_data, update=extra_data)
    session.add(db_user)
    session.commit()
    session.refresh(db_user)
    return db_user

@router.delete("/{user_id}")
async def delete_user(*, session: Session = Depends(get_session), user_id: int):
    user = session.get(User, user_id)
    if not user:
        return JSONResponse(status_code=404, content={"message": "User Not Found"})
    session.delete(user)
    session.commit()
    return{"ok": True}
