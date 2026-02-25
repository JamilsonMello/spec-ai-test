import mongoose, { Schema, Document } from \'mongoose\';
import { ValidationException } from \'../src/exception/validation_exception\';

export interface UserDocument extends Document {
  email: string;
  password?: string;
  name: string;
  role: string;
  passwordRecoveryToken?: string | null;
}

export class UserEntity {
  public email: string;
  public password?: string;
  public name: string;
  public role: string;
  public passwordRecoveryToken?: string | null;

  constructor(email: string, name: string, password?: string, role: string = \'user\', passwordRecoveryToken: string | null = null) {
    UserEntity.validateEmail(email);
    UserEntity.validateName(name);
    if (password) {
      UserEntity.validatePassword(password);
    }

    this.email = email;
    this.name = name;
    this.password = password;
    this.role = role;
    this.passwordRecoveryToken = passwordRecoveryToken;
  }

  private static validateEmail(email: string): void {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      throw new ValidationException(\'Invalid email format\');
    }
  }

  private static validateName(name: string): void {
    if (!name || name.trim().length < 3) {
      throw new ValidationException(\'Name must be a string with at least 3 characters\');
    }
  }

  private static validatePassword(password: string): void {
    if (password.length < 8) {
      throw new ValidationException(\'Password must be at least 8 characters long\');
    }
  }

  // Method to create a UserEntity from plain data, useful for repository
  public static fromObject(data: any): UserEntity {
    return new UserEntity(
      data.email,
      data.name,
      data.password,
      data.role,
      data.passwordRecoveryToken
    );
  }
}

const UserSchema: Schema = new Schema(
  {
    email: {
      type: String,
      required: true,
      unique: true,
      lowercase: true,
      trim: true,
    },
    password: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
      trim: true,
    },
    role: {
      type: String,
      required: true,
      default: \'user\',
      enum: [\'user\', \'admin\'],
    },
    passwordRecoveryToken: {
      type: String,
      default: null,
    },
  },
  {
    timestamps: true,
  }
);

export const User = mongoose.model<UserDocument>(\'User\', UserSchema);
