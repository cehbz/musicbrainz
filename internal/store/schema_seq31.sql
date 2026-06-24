CREATE TABLE alternative_release (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  release INTEGER,
  name TEXT,
  artist_credit INTEGER,
  type INTEGER,
  language INTEGER,
  script INTEGER,
  comment TEXT
);
CREATE TABLE alternative_release_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE alternative_medium (
  id INTEGER,
  medium INTEGER,
  alternative_release INTEGER,
  name TEXT
);
CREATE TABLE alternative_track (
  id INTEGER PRIMARY KEY,
  name TEXT,
  artist_credit INTEGER,
  ref_count INTEGER
);
CREATE TABLE alternative_medium_track (
  alternative_medium INTEGER,
  track INTEGER,
  alternative_track INTEGER
);
CREATE TABLE annotation (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  text TEXT,
  changelog TEXT,
  created TEXT
);
CREATE TABLE application (
  id INTEGER PRIMARY KEY,
  owner INTEGER,
  name TEXT,
  oauth_id TEXT,
  oauth_secret TEXT,
  oauth_redirect_uri TEXT
);
CREATE TABLE area_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE area (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  type INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  ended INTEGER,
  comment TEXT
);
CREATE TABLE area_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE area_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE area_alias (
  id INTEGER PRIMARY KEY,
  area INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE area_annotation (
  area INTEGER,
  annotation INTEGER
);
CREATE TABLE area_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE area_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  area_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE area_attribute (
  id INTEGER PRIMARY KEY,
  area INTEGER,
  area_attribute_type INTEGER,
  area_attribute_type_allowed_value INTEGER,
  area_attribute_text TEXT
);
CREATE TABLE area_containment (
  descendant INTEGER,
  parent INTEGER,
  depth INTEGER
);
CREATE TABLE area_tag (
  area INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE area_tag_raw (
  area INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE artist (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  type INTEGER,
  area INTEGER,
  gender INTEGER,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  ended INTEGER,
  begin_area INTEGER,
  end_area INTEGER,
  discogs_artist_id INTEGER
);
CREATE TABLE artist_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE artist_alias (
  id INTEGER PRIMARY KEY,
  artist INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE artist_annotation (
  artist INTEGER,
  annotation INTEGER
);
CREATE TABLE artist_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE artist_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  artist_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE artist_attribute (
  id INTEGER PRIMARY KEY,
  artist INTEGER,
  artist_attribute_type INTEGER,
  artist_attribute_type_allowed_value INTEGER,
  artist_attribute_text TEXT
);
CREATE TABLE artist_ipi (
  artist INTEGER,
  ipi TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE artist_isni (
  artist INTEGER,
  isni TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE artist_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE artist_tag (
  artist INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE artist_rating_raw (
  artist INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE artist_tag_raw (
  artist INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE artist_credit (
  id INTEGER PRIMARY KEY,
  name TEXT,
  artist_count INTEGER,
  ref_count INTEGER,
  created TEXT,
  edits_pending INTEGER,
  gid TEXT
);
CREATE TABLE artist_credit_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE artist_credit_name (
  artist_credit INTEGER,
  position INTEGER,
  artist INTEGER,
  name TEXT,
  join_phrase TEXT
);
CREATE TABLE artist_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE artist_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE artist_release (
  is_track_artist INTEGER,
  artist INTEGER,
  first_release_date INTEGER,
  catalog_numbers TEXT,
  country_code TEXT,
  barcode INTEGER,
  name TEXT,
  release INTEGER
);
CREATE TABLE artist_release_pending_update (
  release INTEGER
);
CREATE TABLE artist_release_group (
  is_track_artist INTEGER,
  artist INTEGER,
  unofficial INTEGER,
  primary_type_child_order INTEGER,
  primary_type INTEGER,
  secondary_type_child_orders TEXT,
  secondary_types TEXT,
  first_release_date INTEGER,
  name TEXT,
  release_group INTEGER
);
CREATE TABLE artist_release_group_pending_update (
  release_group INTEGER
);
CREATE TABLE autoeditor_election (
  id INTEGER PRIMARY KEY,
  candidate INTEGER,
  proposer INTEGER,
  seconder_1 INTEGER,
  seconder_2 INTEGER,
  status INTEGER,
  yes_votes INTEGER,
  no_votes INTEGER,
  propose_time TEXT,
  open_time TEXT,
  close_time TEXT
);
CREATE TABLE autoeditor_election_vote (
  id INTEGER PRIMARY KEY,
  autoeditor_election INTEGER,
  voter INTEGER,
  vote INTEGER,
  vote_time TEXT
);
CREATE TABLE cdtoc (
  id INTEGER PRIMARY KEY,
  discid TEXT,
  freedb_id TEXT,
  track_count INTEGER,
  leadout_offset INTEGER,
  track_offset TEXT,
  created TEXT
);
CREATE TABLE cdtoc_raw (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  discid TEXT,
  track_count INTEGER,
  leadout_offset INTEGER,
  track_offset TEXT
);
CREATE TABLE country_area (
  area INTEGER PRIMARY KEY
);
CREATE TABLE deleted_entity (
  gid TEXT,
  data TEXT,
  deleted_at TEXT
);
CREATE TABLE edit (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  type INTEGER,
  status INTEGER,
  autoedit INTEGER,
  open_time TEXT,
  close_time TEXT,
  expire_time TEXT,
  language INTEGER,
  quality INTEGER
);
CREATE TABLE edit_data (
  edit INTEGER PRIMARY KEY,
  data TEXT
);
CREATE TABLE edit_note (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  edit INTEGER,
  text TEXT,
  post_time TEXT
);
CREATE TABLE edit_note_change (
  id INTEGER PRIMARY KEY,
  status TEXT,
  edit_note INTEGER,
  change_editor INTEGER,
  change_time TEXT,
  old_note TEXT,
  new_note TEXT,
  reason TEXT
);
CREATE TABLE edit_note_recipient (
  recipient INTEGER,
  edit_note INTEGER
);
CREATE TABLE edit_area (
  edit INTEGER,
  area INTEGER
);
CREATE TABLE edit_artist (
  edit INTEGER,
  artist INTEGER,
  status INTEGER
);
CREATE TABLE edit_event (
  edit INTEGER,
  event INTEGER
);
CREATE TABLE edit_genre (
  edit INTEGER,
  genre INTEGER
);
CREATE TABLE edit_instrument (
  edit INTEGER,
  instrument INTEGER
);
CREATE TABLE edit_label (
  edit INTEGER,
  label INTEGER,
  status INTEGER
);
CREATE TABLE edit_mood (
  edit INTEGER,
  mood INTEGER
);
CREATE TABLE edit_place (
  edit INTEGER,
  place INTEGER
);
CREATE TABLE edit_release (
  edit INTEGER,
  release INTEGER
);
CREATE TABLE edit_release_group (
  edit INTEGER,
  release_group INTEGER
);
CREATE TABLE edit_recording (
  edit INTEGER,
  recording INTEGER
);
CREATE TABLE edit_series (
  edit INTEGER,
  series INTEGER
);
CREATE TABLE edit_work (
  edit INTEGER,
  work INTEGER
);
CREATE TABLE edit_url (
  edit INTEGER,
  url INTEGER
);
CREATE TABLE editor (
  id INTEGER PRIMARY KEY,
  name TEXT,
  privs INTEGER,
  email TEXT,
  website TEXT,
  bio TEXT,
  member_since TEXT,
  email_confirm_date TEXT,
  last_login_date TEXT,
  last_updated TEXT,
  birth_date TEXT,
  gender INTEGER,
  area INTEGER,
  password TEXT,
  ha1 TEXT,
  deleted INTEGER
);
CREATE TABLE old_editor_name (
  name TEXT
);
CREATE TABLE editor_language (
  editor INTEGER,
  language INTEGER,
  fluency TEXT
);
CREATE TABLE editor_preference (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  name TEXT,
  value TEXT
);
CREATE TABLE editor_subscribe_artist (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  artist INTEGER,
  last_edit_sent INTEGER
);
CREATE TABLE editor_subscribe_artist_deleted (
  editor INTEGER,
  gid TEXT,
  deleted_by INTEGER
);
CREATE TABLE editor_subscribe_collection (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  collection INTEGER,
  last_edit_sent INTEGER,
  available INTEGER,
  last_seen_name TEXT
);
CREATE TABLE editor_subscribe_label (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  label INTEGER,
  last_edit_sent INTEGER
);
CREATE TABLE editor_subscribe_label_deleted (
  editor INTEGER,
  gid TEXT,
  deleted_by INTEGER
);
CREATE TABLE editor_subscribe_editor (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  subscribed_editor INTEGER,
  last_edit_sent INTEGER
);
CREATE TABLE editor_subscribe_series (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  series INTEGER,
  last_edit_sent INTEGER
);
CREATE TABLE editor_subscribe_series_deleted (
  editor INTEGER,
  gid TEXT,
  deleted_by INTEGER
);
CREATE TABLE event (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  time TEXT,
  type INTEGER,
  cancelled INTEGER,
  setlist TEXT,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  ended INTEGER
);
CREATE TABLE event_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER,
  event_art_presence TEXT
);
CREATE TABLE event_rating_raw (
  event INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE event_tag_raw (
  event INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE event_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE event_alias (
  id INTEGER PRIMARY KEY,
  event INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE event_annotation (
  event INTEGER,
  annotation INTEGER
);
CREATE TABLE event_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE event_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  event_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE event_attribute (
  id INTEGER PRIMARY KEY,
  event INTEGER,
  event_attribute_type INTEGER,
  event_attribute_type_allowed_value INTEGER,
  event_attribute_text TEXT
);
CREATE TABLE event_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE event_tag (
  event INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE event_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_first_release_date (
  release INTEGER PRIMARY KEY,
  year INTEGER,
  month INTEGER,
  day INTEGER
);
CREATE TABLE recording_first_release_date (
  recording INTEGER PRIMARY KEY,
  year INTEGER,
  month INTEGER,
  day INTEGER
);
CREATE TABLE gender (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE genre (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE genre_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE genre_alias (
  id INTEGER PRIMARY KEY,
  genre INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE genre_annotation (
  genre INTEGER,
  annotation INTEGER
);
CREATE TABLE instrument_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE instrument (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  type INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  comment TEXT,
  description TEXT
);
CREATE TABLE instrument_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE instrument_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE instrument_alias (
  id INTEGER PRIMARY KEY,
  instrument INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE instrument_annotation (
  instrument INTEGER,
  annotation INTEGER
);
CREATE TABLE instrument_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE instrument_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  instrument_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE instrument_attribute (
  id INTEGER PRIMARY KEY,
  instrument INTEGER,
  instrument_attribute_type INTEGER,
  instrument_attribute_type_allowed_value INTEGER,
  instrument_attribute_text TEXT
);
CREATE TABLE instrument_tag (
  instrument INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE instrument_tag_raw (
  instrument INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE iso_3166_1 (
  area INTEGER,
  code TEXT
);
CREATE TABLE iso_3166_2 (
  area INTEGER,
  code TEXT
);
CREATE TABLE iso_3166_3 (
  area INTEGER,
  code TEXT
);
CREATE TABLE isrc (
  id INTEGER PRIMARY KEY,
  recording INTEGER,
  isrc TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE iswc (
  id INTEGER PRIMARY KEY,
  work INTEGER,
  iswc TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE l_area_area (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_artist (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_event (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_genre (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_instrument (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_area_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_artist (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_event (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_genre (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_instrument (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_artist_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_event (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_genre (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_instrument (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_event_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_genre (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_instrument (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_genre_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_instrument (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_label (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_instrument_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_label_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_mood (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_mood_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_place (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_place_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_recording (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_recording_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_release (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_group_release_group (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_group_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_group_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_release_group_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_series_series (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_series_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_series_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_url_url (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_url_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE l_work_work (
  id INTEGER PRIMARY KEY,
  link INTEGER,
  entity0 INTEGER,
  entity1 INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  link_order INTEGER,
  entity0_credit TEXT,
  entity1_credit TEXT
);
CREATE TABLE label (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  label_code INTEGER,
  type INTEGER,
  area INTEGER,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  ended INTEGER,
  discogs_label_id INTEGER
);
CREATE TABLE label_rating_raw (
  label INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE label_tag_raw (
  label INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE label_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE label_alias (
  id INTEGER PRIMARY KEY,
  label INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE label_annotation (
  label INTEGER,
  annotation INTEGER
);
CREATE TABLE label_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE label_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  label_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE label_attribute (
  id INTEGER PRIMARY KEY,
  label INTEGER,
  label_attribute_type INTEGER,
  label_attribute_type_allowed_value INTEGER,
  label_attribute_text TEXT
);
CREATE TABLE label_ipi (
  label INTEGER,
  ipi TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE label_isni (
  label INTEGER,
  isni TEXT,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE label_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE label_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE label_tag (
  label INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE label_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE language (
  id INTEGER PRIMARY KEY,
  iso_code_2t TEXT,
  iso_code_2b TEXT,
  iso_code_1 TEXT,
  name TEXT,
  frequency INTEGER,
  iso_code_3 TEXT
);
CREATE TABLE link (
  id INTEGER PRIMARY KEY,
  link_type INTEGER,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  attribute_count INTEGER,
  created TEXT,
  ended INTEGER
);
CREATE TABLE link_attribute (
  link INTEGER,
  attribute_type INTEGER,
  created TEXT
);
CREATE TABLE link_attribute_type (
  id INTEGER PRIMARY KEY,
  parent INTEGER,
  root INTEGER,
  child_order INTEGER,
  gid TEXT,
  name TEXT,
  description TEXT,
  last_updated TEXT
);
CREATE TABLE link_creditable_attribute_type (
  attribute_type INTEGER PRIMARY KEY
);
CREATE TABLE link_attribute_credit (
  link INTEGER,
  attribute_type INTEGER,
  credited_as TEXT
);
CREATE TABLE link_text_attribute_type (
  attribute_type INTEGER PRIMARY KEY
);
CREATE TABLE link_attribute_text_value (
  link INTEGER,
  attribute_type INTEGER,
  text_value TEXT
);
CREATE TABLE link_type (
  id INTEGER PRIMARY KEY,
  parent INTEGER,
  child_order INTEGER,
  gid TEXT,
  entity_type0 TEXT,
  entity_type1 TEXT,
  name TEXT,
  description TEXT,
  link_phrase TEXT,
  reverse_link_phrase TEXT,
  long_link_phrase TEXT,
  last_updated TEXT,
  is_deprecated INTEGER,
  has_dates INTEGER,
  entity0_cardinality INTEGER,
  entity1_cardinality INTEGER
);
CREATE TABLE link_type_attribute_type (
  link_type INTEGER,
  attribute_type INTEGER,
  min INTEGER,
  max INTEGER,
  last_updated TEXT
);
CREATE TABLE editor_collection (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  editor INTEGER,
  name TEXT,
  public INTEGER,
  description TEXT,
  type INTEGER
);
CREATE TABLE editor_collection_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE editor_collection_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  entity_type TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE editor_collection_collaborator (
  collection INTEGER,
  editor INTEGER
);
CREATE TABLE editor_collection_area (
  collection INTEGER,
  area INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_artist (
  collection INTEGER,
  artist INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_event (
  collection INTEGER,
  event INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_genre (
  collection INTEGER,
  genre INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_instrument (
  collection INTEGER,
  instrument INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_label (
  collection INTEGER,
  label INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_place (
  collection INTEGER,
  place INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_recording (
  collection INTEGER,
  recording INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_release (
  collection INTEGER,
  release INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_release_group (
  collection INTEGER,
  release_group INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_series (
  collection INTEGER,
  series INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_work (
  collection INTEGER,
  work INTEGER,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_collection_deleted_entity (
  collection INTEGER,
  gid TEXT,
  added TEXT,
  position INTEGER,
  comment TEXT
);
CREATE TABLE editor_oauth_token (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  application INTEGER,
  authorization_code TEXT,
  refresh_token TEXT,
  access_token TEXT,
  expire_time TEXT,
  scope INTEGER,
  granted TEXT,
  code_challenge TEXT,
  code_challenge_method TEXT
);
CREATE TABLE medium (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  position INTEGER,
  format INTEGER,
  name TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  track_count INTEGER,
  gid TEXT
);
CREATE TABLE medium_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE medium_attribute_type_allowed_format (
  medium_format INTEGER,
  medium_attribute_type INTEGER
);
CREATE TABLE medium_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  medium_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE medium_attribute_type_allowed_value_allowed_format (
  medium_format INTEGER,
  medium_attribute_type_allowed_value INTEGER
);
CREATE TABLE medium_attribute (
  id INTEGER PRIMARY KEY,
  medium INTEGER,
  medium_attribute_type INTEGER,
  medium_attribute_type_allowed_value INTEGER,
  medium_attribute_text TEXT
);
CREATE TABLE medium_cdtoc (
  id INTEGER PRIMARY KEY,
  medium INTEGER,
  cdtoc INTEGER,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE medium_format (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  year INTEGER,
  has_discids INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE medium_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE mood (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE mood_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE mood_alias (
  id INTEGER PRIMARY KEY,
  mood INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE mood_annotation (
  mood INTEGER,
  annotation INTEGER
);
CREATE TABLE orderable_link_type (
  link_type INTEGER PRIMARY KEY,
  direction INTEGER
);
CREATE TABLE place (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  type INTEGER,
  address TEXT,
  area INTEGER,
  coordinates TEXT,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  ended INTEGER
);
CREATE TABLE place_alias (
  id INTEGER PRIMARY KEY,
  place INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE place_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE place_annotation (
  place INTEGER,
  annotation INTEGER
);
CREATE TABLE place_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE place_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  place_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE place_attribute (
  id INTEGER PRIMARY KEY,
  place INTEGER,
  place_attribute_type INTEGER,
  place_attribute_type_allowed_value INTEGER,
  place_attribute_text TEXT
);
CREATE TABLE place_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE place_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE place_rating_raw (
  place INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE place_tag (
  place INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE place_tag_raw (
  place INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE place_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE replication_control (
  id INTEGER PRIMARY KEY,
  current_schema_sequence INTEGER,
  current_replication_sequence INTEGER,
  last_replication_date TEXT
);
CREATE TABLE recording (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  artist_credit INTEGER,
  length INTEGER,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  video INTEGER
);
CREATE TABLE recording_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE recording_alias (
  id INTEGER PRIMARY KEY,
  recording INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE recording_rating_raw (
  recording INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE recording_tag_raw (
  recording INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE recording_annotation (
  recording INTEGER,
  annotation INTEGER
);
CREATE TABLE recording_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE recording_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  recording_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE recording_attribute (
  id INTEGER PRIMARY KEY,
  recording INTEGER,
  recording_attribute_type INTEGER,
  recording_attribute_type_allowed_value INTEGER,
  recording_attribute_text TEXT
);
CREATE TABLE recording_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE recording_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE recording_tag (
  recording INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE release (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  artist_credit INTEGER,
  release_group INTEGER,
  status INTEGER,
  packaging INTEGER,
  language INTEGER,
  script INTEGER,
  barcode TEXT,
  comment TEXT,
  edits_pending INTEGER,
  quality INTEGER,
  last_updated TEXT,
  discogs_release_id INTEGER
);
CREATE TABLE release_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_alias (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE release_country (
  release INTEGER,
  country INTEGER,
  date_year INTEGER,
  date_month INTEGER,
  date_day INTEGER
);
CREATE TABLE release_unknown_country (
  release INTEGER PRIMARY KEY,
  date_year INTEGER,
  date_month INTEGER,
  date_day INTEGER
);
CREATE TABLE release_raw (
  id INTEGER PRIMARY KEY,
  title TEXT,
  artist TEXT,
  added TEXT,
  last_modified TEXT,
  lookup_count INTEGER,
  modify_count INTEGER,
  source INTEGER,
  barcode TEXT,
  comment TEXT
);
CREATE TABLE release_tag_raw (
  release INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE release_annotation (
  release INTEGER,
  annotation INTEGER
);
CREATE TABLE release_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  release_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_attribute (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  release_attribute_type INTEGER,
  release_attribute_type_allowed_value INTEGER,
  release_attribute_text TEXT
);
CREATE TABLE release_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE release_meta (
  id INTEGER PRIMARY KEY,
  date_added TEXT,
  info_url TEXT,
  amazon_asin TEXT,
  cover_art_presence TEXT
);
CREATE TABLE release_label (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  label INTEGER,
  catalog_number TEXT,
  last_updated TEXT
);
CREATE TABLE release_packaging (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_status (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_tag (
  release INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE release_group (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  artist_credit INTEGER,
  type INTEGER,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  discogs_master_id INTEGER
);
CREATE TABLE release_group_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_group_alias (
  id INTEGER PRIMARY KEY,
  release_group INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE release_group_rating_raw (
  release_group INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE release_group_tag_raw (
  release_group INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE release_group_annotation (
  release_group INTEGER,
  annotation INTEGER
);
CREATE TABLE release_group_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_group_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  release_group_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_group_attribute (
  id INTEGER PRIMARY KEY,
  release_group INTEGER,
  release_group_attribute_type INTEGER,
  release_group_attribute_type_allowed_value INTEGER,
  release_group_attribute_text TEXT
);
CREATE TABLE release_group_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE release_group_meta (
  id INTEGER PRIMARY KEY,
  release_count INTEGER,
  first_release_date_year INTEGER,
  first_release_date_month INTEGER,
  first_release_date_day INTEGER,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE release_group_tag (
  release_group INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE release_group_primary_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_group_secondary_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE release_group_secondary_type_join (
  release_group INTEGER,
  secondary_type INTEGER,
  created TEXT
);
CREATE TABLE script (
  id INTEGER PRIMARY KEY,
  iso_code TEXT,
  iso_number TEXT,
  name TEXT,
  frequency INTEGER
);
CREATE TABLE series (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  comment TEXT,
  type INTEGER,
  ordering_type INTEGER,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE series_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  entity_type TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE series_ordering_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE series_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE series_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE series_alias (
  id INTEGER PRIMARY KEY,
  series INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE series_annotation (
  series INTEGER,
  annotation INTEGER
);
CREATE TABLE series_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE series_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  series_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE series_attribute (
  id INTEGER PRIMARY KEY,
  series INTEGER,
  series_attribute_type INTEGER,
  series_attribute_type_allowed_value INTEGER,
  series_attribute_text TEXT
);
CREATE TABLE series_tag (
  series INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE series_tag_raw (
  series INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE tag (
  id INTEGER PRIMARY KEY,
  name TEXT,
  ref_count INTEGER
);
CREATE TABLE tag_relation (
  tag1 INTEGER,
  tag2 INTEGER,
  weight INTEGER,
  last_updated TEXT
);
CREATE TABLE track (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  recording INTEGER,
  medium INTEGER,
  position INTEGER,
  number TEXT,
  name TEXT,
  artist_credit INTEGER,
  length INTEGER,
  edits_pending INTEGER,
  last_updated TEXT,
  is_data_track INTEGER
);
CREATE TABLE track_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE track_raw (
  id INTEGER PRIMARY KEY,
  release INTEGER,
  title TEXT,
  artist TEXT,
  sequence INTEGER
);
CREATE TABLE medium_index (
  medium INTEGER PRIMARY KEY,
  toc TEXT
);
CREATE TABLE unreferenced_row_log (
  table_name TEXT,
  row_id INTEGER,
  inserted TEXT
);
CREATE TABLE url (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  url TEXT,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE url_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE vote (
  id INTEGER PRIMARY KEY,
  editor INTEGER,
  edit INTEGER,
  vote INTEGER,
  vote_time TEXT,
  superseded INTEGER
);
CREATE TABLE work (
  id INTEGER PRIMARY KEY,
  gid TEXT,
  name TEXT,
  type INTEGER,
  comment TEXT,
  edits_pending INTEGER,
  last_updated TEXT
);
CREATE TABLE work_language (
  work INTEGER,
  language INTEGER,
  edits_pending INTEGER,
  created TEXT
);
CREATE TABLE work_rating_raw (
  work INTEGER,
  editor INTEGER,
  rating INTEGER
);
CREATE TABLE work_tag_raw (
  work INTEGER,
  editor INTEGER,
  tag INTEGER,
  is_upvote INTEGER
);
CREATE TABLE work_alias_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE work_alias (
  id INTEGER PRIMARY KEY,
  work INTEGER,
  name TEXT,
  locale TEXT,
  edits_pending INTEGER,
  last_updated TEXT,
  type INTEGER,
  sort_name TEXT,
  begin_date_year INTEGER,
  begin_date_month INTEGER,
  begin_date_day INTEGER,
  end_date_year INTEGER,
  end_date_month INTEGER,
  end_date_day INTEGER,
  primary_for_locale INTEGER,
  ended INTEGER
);
CREATE TABLE work_annotation (
  work INTEGER,
  annotation INTEGER
);
CREATE TABLE work_gid_redirect (
  gid TEXT,
  new_id INTEGER,
  created TEXT
);
CREATE TABLE work_meta (
  id INTEGER PRIMARY KEY,
  rating INTEGER,
  rating_count INTEGER
);
CREATE TABLE work_tag (
  work INTEGER,
  tag INTEGER,
  count INTEGER,
  last_updated TEXT
);
CREATE TABLE work_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE work_attribute_type (
  id INTEGER PRIMARY KEY,
  name TEXT,
  comment TEXT,
  free_text INTEGER,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE work_attribute_type_allowed_value (
  id INTEGER PRIMARY KEY,
  work_attribute_type INTEGER,
  value TEXT,
  parent INTEGER,
  child_order INTEGER,
  description TEXT,
  gid TEXT
);
CREATE TABLE work_attribute (
  id INTEGER PRIMARY KEY,
  work INTEGER,
  work_attribute_type INTEGER,
  work_attribute_type_allowed_value INTEGER,
  work_attribute_text TEXT
);
