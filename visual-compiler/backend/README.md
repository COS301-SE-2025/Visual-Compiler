## **Steps to run the API**
1. Change Directory into the backend folder `cd visual-compiler/backend`.
2. Inside the backend create a `.env` file.
3. Write this code into the .env file:
```env
Mongo_username=<student number>
Mongo_password=<password when connecting to the cluster>
Mongo_URI=mongodb+srv://%s:%s@visualcompiler.ip4jaot.mongodb.net/?retryWrites=true&w=majority&appName=VisualCompiler
```
4. Now in the terminal, still inside the backend folder, type the command: `air`.
5. The API should be running at this point.