import express, { Request, Response } from 'express';
import { loginHandler } from '../../gateway/handler/login_handler';
import { passwordRecoveryHandler } from '../../gateway/handler/password_recovery_handler';

const router = express.Router();

router.post('/login', loginHandler);
router.post('/reset-password', passwordResetHandler);
router.post('/users', registerHandler);

export { router as authRouter };
