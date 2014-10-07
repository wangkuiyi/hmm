# HMM Trainer Users Guide

Author: Yi Wang

Last Update: 09/25/2014

## Preamble

THIS PROGRAM USES SOME INTELLECTURAL PROPERTIES INVENTED DURING MY PHD
RESEARCH.  IT IS OK FOR US TO USE IT IN OUR RESEARCH AND EXPERIMENTS,
BUT DO NOT DISTRIBUTE IT.  BEFORE WE MOVE TO PROUDCTION, WE MIGHT NEED
TO WRITE A NEW TRAINING SYSTEM OR REQUEST IP FROM ITS OWNER (WHICH
MIGHT NOT BE ME, BUT TSINGHUA UNIVERSITY).

## Training Corpus

The training corpus is a JSON file of a sequence of `core.Instance`
defined in `core/instance.go`:

    type Instance struct {
	    Obs     [][]Observed // A sequence of observed.
    	Periods []int        // Periods of above observed.
	    index   []int        // Map time t to an observed.
    }

    type Observed map[string]int

where each instance represent a time-series data, e.g., the work
experience of a LinkedIn member.  Consider a member has three work
experience which last 3, 4, and 1 years respectively, the `Period`
field would be:

    "Periods": [3, 4, 1]

If in the first period, we want to consider *tokens in company name*
and *tokens in title* as output observables, we would say we have two
*channels*.  Consider that the above member worked at Google as a
research techlead manager, at Tencent as Engineering Director, and at
LinkedIn as Senior Staff Business Analyst, then the `Obs` field would
be:

    "Obs": [
      [
        {"google":1},
        {"research":1, "techlead":1, "manager":1}
      ],
      [
        {"tencent":1},
        {"engineering":1, "director":1}
      ],
      [
        {"linkedin":1},
        {"senior":1, "staff":1, "business":1, "analyst":1}
      ]
    ]

An example corpus file is at `trainer/testdata/corpus.json`:

     [
        {
            "Obs": [
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ],
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ],
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ]
            ],
            "Periods": [1,1,1,1,1,1]
        },
        {
            "Obs": [
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ],
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ],
                [
                    {"apple": 1}
                ],
                [
                    {"orange": 1}
                ]
            ],
            "Periods": [1,1,1,1,1,1]
        }
    ]

In this example, there are two instances, which are identical.  All
observables of these instances have only one channel, which
corresponds to the name of fruits.

## Use Training System

The current training system runs on a single computer (but not a
cluster of computers).   The supported command flags include:

     -addr=":6060": Listening address
     -corpus="": Corpus file in JSON format
     -iter=20: Number of EM iterations
     -model="": Model file in JSON format
     -states=2: Number of hidden states

The following example command line launches a training job that learns
two hidden states (two clusters) from above example corpus using 100
EM iterations, and writes learned model into `/tmp/model`.

    $GOPATH/bin/trainer -corpus=testdata/corpus.json \
                        -model=/tmp/model \
                        -states=2 \
                        -iter=100 \

If you do not specify `-model`, the model will be printed to standard
output (i.e., your screen in most cases).

## Interpret the Model

The training system outputs model in JSON format. For example, the
model generated from above training jobs is like:

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

where

1. `Σγo` represents the learned output probability distributions.  In
   above example, there are two states, so two output distributions --
   the first concentrate on `orange` and the other one concentrate on
   `orange`.

1. `Σξ` represents the number of transitions.  In above example, we
   see that it is impossible to transit from state 0 to 0, or from 1
   to 1; but possible to transit from 0 to 1 with probability
   4/(4+6)=40%, and from 1 to 0 with probability 6/(4+6)=60%.  This
   matches the fact in our training corpus that transitions are always
   from "apple" ro "orange", or from "orange" to "apple", but no
   transitions from "apple" to "apple", or from "orange" to "orange".

1. `Σγ`represents the the number of times that transitions are from a
   certain hidden state.  In above example, 4 transitons start from
   state 0, and 6 transitons start from 1.

1. `S1` and `S1Sum` count the number of times that an instance start
   from a certain hidden state.  In above example, both instances
   start from hidden state 1, which, according to `Σγo`, is about
   "apple".  This fact again matches our training corpus.
