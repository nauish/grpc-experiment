import express from 'express';
const app = express();
app.get('/', async (req, res) => {
  const arr = [];
  const promises = [];
  for (let index = 0; index < 10; index++) {
    const promise = fetch('http://localhost:8080')
      .then((response) => response.json())
      .then((data) => arr.push(data));
    promises.push(promise);
  }
  await Promise.all(promises);
  res.status(200).send(arr);
});
app.listen(5002, () => console.log(`Server running on port 5002`));
