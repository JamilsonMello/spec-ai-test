const User = require('../models/user.model');

const listUsers = async (req, res) => {
  try {
    const { name, email } = req.query;
    const page = parseInt(req.query.page) || 1;
    const limit = parseInt(req.query.limit) || 30; // Max 30 items per page as per spec
    const skip = (page - 1) * limit;

    let query = {};
    if (name) {
      query.name = { $regex: name, $options: 'i' }; // Case-insensitive search
    }
    if (email) {
      query.email = { $regex: email, $options: 'i' }; // Case-insensitive search
    }

    const users = await User.find(query)
      .sort({ createdAt: -1 }) // Default sorting by created_at descending
      .skip(skip)
      .limit(limit)
      .select('-password -recoveryToken'); // Omit sensitive fields

    const totalUsers = await User.countDocuments(query);
    const totalPages = Math.ceil(totalUsers / limit);

    res.status(200).json({
      users,
      currentPage: page,
      totalPages,
      totalUsers,
    });
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

module.exports = {
  listUsers,
};
