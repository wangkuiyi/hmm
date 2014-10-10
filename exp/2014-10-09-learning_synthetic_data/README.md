
I tested learning from synthetic data to see if our HMM training
program can reach global optimal estimate.

A very simple model (with only two hidden states) used to
generate/synthesize training corpus is as follows:

    {
      "S1": [
        0,
        2
      ],
      "S1Sum": 2,
      "Σγ": [
        4,
        6
      ],
      "Σξ": [
        [
          0,
          4
        ],
        [
          6,
          0
        ]
      ],
      "Σγo": [
        [
          {
            "Hist": {
              "orange": 6
            },
            "Sum": 6
          }
        ],
        [
          {
            "Hist": {
              "apple": 6
            },
            "Sum": 6
          }
        ]
      ]
    }

I synthesized 20 instances, each with length 6 (years), and
multinomial sample cardinality is 1.  Then I learned a model, shows as
follows:

    {
      "S1": [
        20,
        0
      ],
      "S1Sum": 20,
      "Σγ": [
        60,
        40
      ],
      "Σξ": [
        [
          0,
          60
        ],
        [
          40,
          0
        ]
      ],
      "Σγo": [
        [
          {
            "Hist": {
              "apple": 60
            },
            "Sum": 60
          }
        ],
        [
          {
            "Hist": {
              "orange": 60
            },
            "Sum": 60
          }
        ]
      ]
    }

### Conclusion

The learned model matches exactly with the ground-truth model.  Counts
in learned model is ten times of those in the ground-truth model,
because the synthetic data is ten times bigger than the data used to
estimate the ground-truth model.

In ground-truth model, the first state represents orange and the
second represents apple; wheresa in the learned model, the first state
represents apple and the second represents orange.  This is reasonable
and completely acceptable, as we are doing clustering and we have no
way to gaurrentee the order of clusters.
