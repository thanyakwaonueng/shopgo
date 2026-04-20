TODAY WORK: since I think these project is mostly sufficient now so I will try to do the frontend by my self using shopgo-frontend from mr.saza as example

NOTE: I know that .env should be commited as .env.example

NOTE: service code with db transaction is not in purely repository pattern yet(it is for most part but there exists transaction call in the service layer, hence I call it not pure-repository), this happen because I want some business logic to be atomic(either all passed or all roll backed)

NOTE: I left out unit test of login, register, getMe and refreshToken service for now since the code is literally just mapping a request to db with minimal(if any) business logic at the moment
