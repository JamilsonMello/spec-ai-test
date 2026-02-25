import { Request, Response, NextFunction } from 'express';
import { UserDocument } from '../../entity/user';

declare module 'express-serve-static-core' {
  interface Request {
    user?: UserDocument; // Or a more specific user type from your auth system
  }
}

export const adminMiddleware = (req: Request, res: Response, next: NextFunction) => {
  // Assuming user information is attached to req.user by a previous authentication middleware
  if (!req.user || req.user.role !== 'admin') {
    return res.status(403).json({ error: 'Access denied. Administrator privileges required.' });
  }
  next();
};
