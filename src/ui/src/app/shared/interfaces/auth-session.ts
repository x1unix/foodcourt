/**
 * User auth session
 */
export interface IAuthSession {
  /**
   * Session time to live (UNIX UTC timestamp)
   */
  ttl: number;

  /**
   * Token ID
   */
  token: string;

  /**
   * Is authorized (currently always true)
   */
  authorized: boolean;

  /**
   * Current user ID
   */
  userId: number;
}
