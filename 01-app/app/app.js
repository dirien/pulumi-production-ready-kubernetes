const express = require('express');
const bodyParser = require('body-parser');

// Create Express server
const app = express();

// Middleware
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));

// GET route
app.get('/', (req, res) => {
    res.status(200).send('Hello Pulumi World!');
});

// POST route
app.post('/', (req, res) => {
    const {message} = req.body;
    if (!message) {
        return res.status(400).send('Message is required');
    }
    res.status(200).send(`Received message: ${message}`);
});

// Start server
app.listen(3000, () => {
    console.log('Server started on port 3000');
});
