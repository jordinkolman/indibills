from fastapi import APIRouter, Depends, Query
from fastapi.responses import JSONResponse
from sqlalchemy.exc import NoResultFound
from sqlmodel import Session, select

from ..dependencies import get_session
from ..models import (
    ManualAccount,
    ManualTransactionCreate,
    ManualTransaction,
    ManualTransactionUpdate,
    Message,
)


router = APIRouter(
    prefix="/manual_transactions",
    tags=["manual_transactions"],
    responses={404: {"model": Message}},
)


@router.get("/account/{account_id}", response_model=list[ManualTransaction])
async def get_manual_transactions_by_account(
    *,
    session: Session = Depends(get_session),
    offset: int = 0,
    limit: int = Query(default=25, le=100),
    account_id: int,
):
    account = session.get(ManualAccount, account_id)
    if not account:
        return JSONResponse(status_code=404, content={"message": "Account Not Found"})
    transactions = session.exec(
        select(ManualTransaction)
        .where(ManualTransaction.account_id == account_id)
        .offset(offset)
        .limit(limit)
    ).all()
    if not transactions:
        return JSONResponse(
            status_code=404,
            content={"message": "No Transactions Found for this Account"},
        )
    return transactions


@router.post("/", response_model=ManualTransaction)
async def create_manual_transaction(
    *, session: Session = Depends(get_session), transaction: ManualTransactionCreate
):
    db_transaction = ManualTransaction.model_validate(transaction)
    session.add(db_transaction)
    session.commit()
    session.refresh(db_transaction)
    return db_transaction


@router.get("/{transaction_id}", response_model=ManualTransaction)
async def get_manual_transaction_by_id(
    *, session: Session = Depends(get_session), transaction_id: int
):
    try:
        transaction = session.exec(
            select(ManualTransaction).where(ManualTransaction.id == transaction_id)
        ).one()
    except NoResultFound:
        return JSONResponse(
            status_code=404, content={"message": "Transaction Not Found"}
        )
    return transaction


@router.patch("/{transaction_id}", response_model=ManualTransaction)
async def update_manual_transaction(
    *,
    session: Session = Depends(get_session),
    transaction_id: int,
    transaction: ManualTransactionUpdate,
):
    db_transaction = session.get(ManualTransaction, transaction_id)
    if not db_transaction:
        return JSONResponse(
            status_code=404, content={"message": "Transaction Not Found"}
        )
    transaction_data = transaction.model_dump(exclude_unset=True)
    db_transaction.sqlmodel_update(transaction_data)
    session.add(db_transaction)
    session.commit()
    session.refresh(db_transaction)
    return db_transaction


@router.delete("/{transaction_id}")
async def delete_manual_transaction(
    *, session: Session = Depends(get_session), transaction_id: int
):
    transaction = session.get(ManualTransaction, transaction_id)
    if not transaction:
        return JSONResponse(
            status_code=404, content={"message": "Transaction Not Found"}
        )
    session.delete(transaction)
    session.commit()
    return {"ok": True}
