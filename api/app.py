from contextlib import asynccontextmanager

from sqlmodel import Session, SQLModel, select
from fastapi import Depends, FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from passlib.hash import pbkdf2_sha256

import models
from engine import engine

def create_db_and_tables():
    SQLModel.metadata.create_all(engine)

@asynccontextmanager
async def lifespan(app: FastAPI):
    create_db_and_tables()
    yield

app = FastAPI(lifespan=lifespan)

origins = [
    'http://127.0.0.1:5500',
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
    expose_headers=["*"],
)

def get_session():
    with Session(engine) as session:
        yield session

# User Queries
@app.post("/users/", response_model=models.UserPublic)
def create_user(*, session: Session = Depends(get_session), user: models.UserCreate):
    hashed_password = pbkdf2_sha256.hash(user.password)
    extra_data = {"hashed_password": hashed_password}
    db_user = models.User.model_validate(user, update=extra_data)
    session.add(db_user)
    session.commit()
    session.refresh(db_user)
    return db_user

@app.get("/users/")
def get_users(session: Session = Depends(get_session), offset: int = 0, limit: int = Query(default=100, le=100)):
    users = session.exec(select(models.User).offset(offset).limit(limit)).all()
    results = []
    for user in users:
        new_user = models.UserPublic(
            id = user.id,
            username = user.username,
            email = user.email,
            created_at=user.created_at,
            updated_at=user.updated_at,
            items = user.items,
            assets = user.assets,
            manual_accounts=user.manual_accounts,
        )
        results.append(new_user)

    return {"users": results}

@app.get("/users/{username}", response_model=models.UserPublic)
def get_user_by_username(*, session: Session = Depends(get_session), username: str):
    return session.exec(select(models.User).where(models.User.username == username)).one()

@app.get("/users/id/{user_id}", response_model=models.UserPublic)
def get_user_by_id(*, session: Session = Depends(get_session), user_id: int):
    return session.exec(select(models.User).where(models.User.id == user_id)).one()

@app.get("/users/email/{user_email}", response_model=models.UserPublic)
def get_user_by_email(*, session: Session = Depends(get_session), user_email: str):
    return session.exec(select(models.User).where(models.User.email == user_email)).one()

@app.patch("/users/{user_id}", response_model=models.UserPublic)
def update_user(*, session: Session = Depends(get_session), user_id: int, user: models.UserUpdate):
    db_user = session.get(models.User, user_id)
    if not db_user:
        return HTTPException(status_code=404, detail="User Not Found")
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

@app.delete("/users/{user_id}")
def delete_user(*, session: Session = Depends(get_session), user_id: int):
    user = session.get(models.User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User Not Found")
    session.delete(user)
    session.commit()
    return{"ok": True}
