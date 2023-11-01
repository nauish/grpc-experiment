import express from 'express';

const app = express();
app.get('/', (req, res) => res.json({ message: 'Hello node!' }));
app.listen(8080, () => console.log(`Server running on port 8080`));
