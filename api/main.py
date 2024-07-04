from fastapi import FastAPI

from .routers import users, manual_accounts, households, manual_transactions

from .dependencies import lifespan


app = FastAPI(lifespan=lifespan)

app.include_router(users.router)
app.include_router(manual_accounts.router)
app.include_router(households.router)
app.include_router(manual_transactions.router)

@app.get("/")
async def root():
    return {"message": "Hello World"}
