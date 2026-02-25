import express, { Request, Response } from 'express';
import { User } from '../../entity/user';


const router = express.Router();

router.get('/users', async (req: Request, res: Response) => {
  try {
    const { name, email, page = 1 } = req.query;

    // Check if user is admin
    const user = await User.findById(req.user.id);
    if (!user || user.role !== 'admin') {
      return res.status(403).send('Access denied');
    }

    const filter: any = {};
    if (name) {
      filter.name = { $regex: name, $options: 'i' };
    }
    if (email) {
      filter.email = { $regex: email, $options: 'i' };
    }

    const limit = 30;
    const skip = (Number(page) - 1) * limit;

    const users = await User.find(filter)
      .sort({ createdAt: -1 })
      .skip(skip)
      .limit(limit)
      .select('-password -passwordRecoveryToken');

    res.send(users);
  } catch (error) {
    console.error(error);
    res.status(500).send('Internal Server Error');
  }
});

export { router as userRouter };