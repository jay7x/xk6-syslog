// index.d.ts

/**
 * Supported transport protocols for connecting to a syslog server.
 */
type Transport = "udp" | "tcp" | "tls";

/**
 * TLS configuration options for secure connections.
 */
interface TLSConfig {
  /**
   * PEM-encoded CA certificate (optional).
   * If provided, it will be used as the root of trust.
   */
  ca?: string;

  /**
   * PEM-encoded client certificate (optional).
   * Required if using mutual TLS authentication.
   */
  cert?: string;

  /**
   * PEM-encoded private key for the client certificate (optional).
   * Required if using mutual TLS authentication.
   */
  key?: string;

  /**
   * Server name used for SNI and peer certificate validation (optional).
   */
  serverName?: string;

  /**
   * If true, skips TLS certificate verification (insecure, use with caution).
   */
  insecureSkipVerify?: boolean;
}

/**
 * Configuration object for connecting to a syslog server.
 */
interface SyslogConfig {
  /**
   * Transport protocol to use: "udp", "tcp", or "tls".
   */
  transport: Transport;

  /**
   * Connection timeout in seconds (optional).
   */
  timeout?: number;

  /**
   * TLS configuration (only used if transport is "tls").
   */
  tls?: TLSConfig;
}

/**
 * Represents an open connection to a syslog server.
 */
interface SyslogConnection {
  /**
   * Sends raw binary data to the syslog server.
   * @param data - The data to send as a Uint8Array or Buffer-like array.
   */
  send(data: ArrayLike<number>): void;

  /**
   * Closes the connection gracefully.
   */
  close(): void;
}

/**
 * Module namespace for interacting with syslog servers.
 */
declare namespace syslog {
  /**
   * Establishes a new connection to a syslog server.
   * @param address - The address of the syslog server (e.g., "localhost:514").
   * @param config - Configuration options for the connection.
   * @returns A {@link SyslogConnection} object for sending and closing.
   */
  function connect(address: string, config: SyslogConfig): SyslogConnection;
}

export = syslog;