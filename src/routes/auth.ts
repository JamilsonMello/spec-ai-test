import express, { Request, Response } from 'express';
import { loginHandler } from '../../gateway/handler/login_handler';
import { passwordRecoveryHandler } from '../../gateway/handler/password_recovery_handler';
import { passwordResetHandler } from '../../gateway/handler/password_reset_handler';
import { registerHandler } from '../../gateway/handler/register_handler';
import { listUsersHandler } from '../../gateway/handler/list_users_handler';
import { adminMiddleware } from '../../gateway/handler/admin_middleware';

// Assuming there's an authentication middleware that populates req.user
const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
  // Placeholder for actual authentication logic
  // In a real application, this would validate a token and attach user info to req.user
  // For this spec, we'll mock an authenticated user with a role
  if (req.headers.authorization === 'Bearer admin_token') {
    req.user = { role: 'admin' }; // Mock admin user
  } else if (req.headers.authorization === 'Bearer user_token') {
    req.user = { role: 'user' }; // Mock regular user
  }
  next();
};

const router = express.Router();

router.post('/login', loginHandler);
router.post('/reset-password', passwordResetHandler);
router.post('/users', registerHandler);
router.get('/users', authMiddleware, adminMiddleware, listUsersHandler);

export { router as authRouter };
