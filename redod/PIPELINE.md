# Redod Pipeiline

On start-up we already have one connection, and throughout the run, we
will have more as the sub-shells call redo themselves. Each connection
delivers a series of requests each containing a target that need to be
updated (or similar). Each request receives exactly one response.

The first step when a target is received is to demand it. Demanding
the target brings its dependencies in to the current generation. All
the newly demanded targets start in the "unknown" state.

The second step is scanning. Scanning a target is to stat the file and
compare it to the old stat information. This also lets us know if the
target was created (if the old stat was missing). If a scanned target
is changed, all of its dependents are marked for update. Scanning
always proceeds by scanning the targets with no un-scanned
dependencies, this ensures the minimum amount of scanning. After a
target is scanned, it is considered up-to-date.

The third step is running. Running means executing the do-file for the
target. Running always proceeds from targets marked for update with
all the dependencies marked up-to-date. After a target runs
successfully, it is considered up-to-date. On failure, the target is
considered broken.

The final step is responding. Once a target is marked up-to-date or
broken, any requests on that target need a response. Connections with
an outstanding request for the target are sent a status message for
the target.

Each step of the pipeline (except the first and last) are goroutines
that use SQLite to find the next targets for their operation. They are
kicked via a channel so that they block when there is no work for them
to do.

The first and last stages are simply a map from targets to lists of
channels. When the connection sends its request, the targets are
demanded immediately, then the connection is placed in the map once
for each target. On the last stage when the fate of the target is
finally known, the connections waiting on that target are all sent a
message with the resolution.
