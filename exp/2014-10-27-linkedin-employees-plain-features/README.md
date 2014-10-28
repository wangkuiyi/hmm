This experiment contains a program `generate_training_data` that
converts `data/pos_and_edu_for_LI_employees v2.txt` into
`./corpus.json`.


We use `plot` to visualize the model, `./model.json`, trained from
`./corpus.json`. The result is `model.pdf`.

It is notable that as we use properties of
`data/pos_and_edu_for_LI_employees v2.txt` as plain features, and all
members in this file are LinkedIn employees, so the feature
`company=linkedin` dominates the clustering result, as shown by
`model.pdf`.
