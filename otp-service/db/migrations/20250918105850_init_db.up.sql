CREATE TABLE session (
  id              VARCHAR(64) PRIMARY KEY,
  -- แนะนำตัด human_verified หรือ status ที่ซ้ำกันออกสักอัน
  human_verified  TINYINT(1) NOT NULL DEFAULT 0,
  idcard_hash     VARBINARY(32) NULL,
  phone_hash      VARBINARY(32) NULL,
  phone_verified  TINYINT(1) NOT NULL DEFAULT 0,
  email_hash      VARBINARY(32) NULL,
  email_verified  TINYINT(1) NOT NULL DEFAULT 0,
  status          ENUM('human_verified','user_verified','sms_otp_verified','document_requested','completed')
                  NOT NULL DEFAULT 'human_verified',
  created_at      DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  expires_at      DATETIME(6) NOT NULL,
  updated_at      DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE otp_code (
  id             BIGINT PRIMARY KEY AUTO_INCREMENT,
  session_id     VARCHAR(64) NOT NULL,
  channel        ENUM('sms','email') NOT NULL,
  otp_ref        VARCHAR(64) NOT NULL,
  otp_code_hash  VARBINARY(32) NOT NULL,
  max_attempts   INT NOT NULL DEFAULT 6,
  attempts       INT NOT NULL DEFAULT 0,
  expires_at     DATETIME(6) NOT NULL,
  created_at     DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at     DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  CONSTRAINT fk_otp_session
    FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,

  UNIQUE KEY idx_session_id_channel_otp_ref (session_id, channel, otp_ref)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;