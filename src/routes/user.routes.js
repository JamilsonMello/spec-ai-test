const express = require('express');
const router = express.Router();
const userController = require('../controllers/user.controller');
const authMiddleware = require('../middlewares/auth.middleware');

// Route for listing users (requires admin privileges)
router.get('/users', authMiddleware, userController.listUsers);

module.exports = router;
