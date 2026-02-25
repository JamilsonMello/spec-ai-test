const jwt = require('jsonwebtoken');

const auth = (req, res, next) => {
  // In a real application, you would extract the token from the request header
  // For this spec, we'll assume req.user is already populated by a previous auth middleware
  // and contains user information including the role.

  if (!req.user || req.user.role !== 'admin') {
    return res.status(403).json({ message: 'Access denied. Administrator privileges required.' });
  }
  next();
};

module.exports = auth;
