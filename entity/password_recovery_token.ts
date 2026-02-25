import mongoose, { Schema, Document } from 'mongoose';

export interface PasswordRecoveryTokenDocument extends Document {
  token: string;
  userId: mongoose.Schema.Types.ObjectId;
  createdAt: Date;
}

const PasswordRecoveryTokenSchema: Schema = new Schema(
  {
    token: {
      type: String,
      required: true,
      unique: true,
    },
    userId: {
      type: mongoose.Schema.Types.ObjectId,
      required: true,
      ref: 'User',
    },
    createdAt: {
      type: Date,
      required: true,
      default: Date.now,
      expires: '24h', // TTL index for 24 hours
    },
  },
  {
    timestamps: true,
  }
);

export const PasswordRecoveryToken = mongoose.model<PasswordRecoveryTokenDocument>(
  'PasswordRecoveryToken',
  PasswordRecoveryTokenSchema
);

export class PasswordRecoveryTokenEntity {
  private token: string;
  private userId: mongoose.Schema.Types.ObjectId;
  private createdAt: Date;

  constructor(token: string, userId: mongoose.Schema.Types.ObjectId, createdAt: Date) {
    this.token = token;
    this.userId = userId;
    this.createdAt = createdAt;
  }

  public isExpired(): boolean {
    const twentyFour Hours = 24 * 60 * 60 * 1000; // 24 hours in milliseconds
    return (new Date().getTime() - this.createdAt.getTime()) > twentyFourHours;
  }

  public getToken(): string {
    return this.token;
  }

  public getUserId(): mongoose.Schema.Types.ObjectId {
    return this.userId;
  }

  public getCreatedAt(): Date {
    return this.createdAt;
  }
}
