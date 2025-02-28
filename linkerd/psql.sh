#!/bin/bash
function finish {
    sync_unlock.sh
}
if [ -z "$TRAP" ]
then
  sync_lock.sh || exit -1
  trap finish EXIT
  export TRAP=1
fi
set -o pipefail
> errors.txt
> run.log
GHA2DB_PROJECT=linkerd IDB_DB=linkerd PG_DB=linkerd GHA2DB_LOCAL=1 ./structure 2>>errors.txt | tee -a run.log || exit 1
GHA2DB_PROJECT=linkerd IDB_DB=linkerd PG_DB=linkerd GHA2DB_LOCAL=1 ./gha2db 2017-01-23 0 today now 'linkerd' 2>>errors.txt | tee -a run.log || exit 2
GHA2DB_PROJECT=linkerd IDB_DB=linkerd PG_DB=linkerd GHA2DB_LOCAL=1 GHA2DB_EXACT=1 ./gha2db 2016-01-13 0 2017-01-24 0 'BuoyantIO/linkerd' 2>>errors.txt | tee -a run.log || exit 3
GHA2DB_PROJECT=linkerd IDB_DB=linkerd PG_DB=linkerd GHA2DB_LOCAL=1 GHA2DB_MGETC=y GHA2DB_SKIPTABLE=1 GHA2DB_INDEX=1 ./structure 2>>errors.txt | tee -a run.log || exit 4
GHA2DB_PROJECT=linkerd PG_DB=linkerd ./shared/setup_repo_groups.sh 2>>errors.txt | tee -a run.log || exit 5
GHA2DB_PROJECT=linkerd IDB_DB=linkerd PG_DB=linkerd ./shared/import_affs.sh 2>>errors.txt | tee -a run.log || exit 6
GHA2DB_PROJECT=linkerd PG_DB=linkerd ./shared/setup_scripts.sh 2>>errors.txt | tee -a run.log || exit 7
GHA2DB_PROJECT=linkerd PG_DB=linkerd ./shared/get_repos.sh 2>>errors.txt | tee -a run.log || exit 8
GHA2DB_PROJECT=linkerd PG_DB=linkerd GHA2DB_LOCAL=1 ./pdb_vars || exit 9
./devel/ro_user_grants.sh linkerd || exit 10
./devel/psql_user_grants.sh devstats_team linkerd || exit 11
