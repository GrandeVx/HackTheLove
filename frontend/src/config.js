const googleClientId = import.meta.env.GOOGLE_CLIENT_ID ||
  '443648413060-db3ivje6uto4h1jf0f11e13hb4opmhep.apps.googleusercontent.com';

const timeReleaseMatch = import.meta.env.TIME_RELEASE_MATCH
  ? new Date(import.meta.env.TIME_RELEASE_MATCH)
  : new Date("2025-02-12T21:00:00Z"); // 2025-02-14T00:00:00Z

export { googleClientId, timeReleaseMatch };
