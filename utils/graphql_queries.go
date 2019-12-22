package utils

const UpsertInClusterTable string = `mutation {
										upsert: upsert_clusters  (docs : $docs) @postgres {status error}
																	}`

