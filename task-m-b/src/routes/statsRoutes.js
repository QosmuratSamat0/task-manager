const express = require('express');
const { getStats } = require('../controllers/statsController');
const { authenticate } = require('../middleware/auth');

const router = express.Router();

router.get('/', authenticate, getStats);

module.exports = router;
