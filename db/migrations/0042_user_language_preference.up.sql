-- Per-user language preference for the dashboard (English default, Amharic supported).
alter table users
  add column if not exists language_preference text not null default 'en'
  check (language_preference in ('en','am'));
