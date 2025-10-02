# Database Management - Simple Guide

**TL;DR: Just run `make dev` and the database sets itself up automatically.**

## Common Tasks

### Starting Development
```bash
make dev
```
Database automatically creates all tables, indexes, and constraints on first run.

### Backup Database (before deploying changes)
```bash
make db-backup
```
Creates timestamped backup in `./backups/`

### Restore from Backup
```bash
make db-restore
# Enter filename when prompted
```

### View/Query Database
```bash
make db-shell
# Now run SQL: SELECT * FROM users;
```

### Something Broke - Start Fresh
```bash
docker-compose down -v  # Delete everything
make dev                # Start fresh with clean database
```

## Files You Should Know About

- **docker-entrypoint-initdb.sh** - Runs automatically on first startup, creates all tables
- **backup.sh** - Used by `make db-backup`
- **restore.sh** - Used by `make db-restore`

## Files You Can Ignore

The `migrations/` folder and old SQL files are legacy - the initialization script handles everything now.

## When Do I Need Migrations?

**Short answer:** Probably never for development.

**Long answer:** Only if you're changing the schema on a live production database with real user data. For now, just modify `docker-entrypoint-initdb.sh` and restart fresh.

## Production Notes

Before deploying schema changes to production:

1. **Backup first:** `make db-backup`
2. **Test locally:** `docker-compose down -v && make dev`
3. **Deploy changes**
4. **If something breaks:** `make db-restore`

---

For detailed database schema documentation, see [DATABASE.md](../DATABASE.md)
