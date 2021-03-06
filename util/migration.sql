DROP TABLE IF EXISTS event_invitees;
DROP TABLE IF EXISTS event_notifications;
DROP TABLE IF EXISTS event_media;
DROP TABLE IF EXISTS events;

CREATE TABLE events(
  event_id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
  slug VARCHAR(255) NOT NULL,
  UNIQUE(slug),
  
  invitation_message TEXT,

  creator_uid VARCHAR(255) NOT NULL,
  
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE event_invitees(
	invitee_id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
	name VARCHAR(255) NOT NULL,
	slug VARCHAR(255) NOT NULL,
	
	event_id BINARY(16),
	FOREIGN KEY (event_id) REFERENCES events(event_id),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE event_notifications(
  notification_id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
  type VARCHAR(64) NOT NULL,
  data JSON,
  
  event_id BINARY(16),
	FOREIGN KEY (event_id) REFERENCES events(event_id),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE event_media(
  media_id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),

  path TEXT,
  video_poster_path TEXT,
  type VARCHAR(64) NOT NULL,
  
  event_id BINARY(16),
	FOREIGN KEY (event_id) REFERENCES events(event_id),

  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);