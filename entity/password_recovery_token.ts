import { ValidationException } from \"../src/exception/validation_exception\";

export class PasswordRecoveryTokenEntity {
  public token: string;
  public createdAt: Date;
  private static readonly EXPIRATION_HOURS = 24;

  constructor(token: string, createdAt: Date) {
    if (!token || token.trim().length === 0) {
      throw new ValidationException(\"Token cannot be empty\");
    }
    this.token = token;
    this.createdAt = createdAt;
  }

  public isExpired(): boolean {
    const now = new Date();
    const expirationTime = new Date(this.createdAt);
    expirationTime.setHours(expirationTime.getHours() + PasswordRecoveryTokenEntity.EXPIRATION_HOURS);
    return now > expirationTime;
  }

  public static fromObject(data: any): PasswordRecoveryTokenEntity {
    return new PasswordRecoveryTokenEntity(data.token, data.createdAt);
  }
}
