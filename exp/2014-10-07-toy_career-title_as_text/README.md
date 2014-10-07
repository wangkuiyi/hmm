

I crafted a toy career data set with two persons both working in the
`internet` industry, but one follows the technical path, whereas the
other one changes to a manger after his techincal path reaches `senior
software engineer`.

    [
	{
	    "Obs": [
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,   // I believe seniority-value as features would
			"software": 1    // makes the learning easiler, but I do not 
		    }                    // use it; instead, I use text features for toy.
		],
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,
			"senior": 1,
			"software": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,
			"software": 1,
			"staff": 1
		    }
		],
		[
		    {
			"internet": 1   // This guy keeps in technical track and grows
		    },                  // from a engineer to senior staff.
		    {
			"engineer": 1,
			"senior": 1,
			"software": 1,
			"staff": 1
		    }
		]
	    ],
	    "Periods": [                 // This s a four-year career path.
		1,
		1,
		1,
		1
	    ]
	},
	{
	    "Obs": [
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,
			"software": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,
			"senior": 1,
			"software": 1
		    }
		],
		[
		    {
			"internet": 1    // This guy changes from technical track to
		    },                   // manager track after he's promoted to 
		    {                    // senior engineer.
			"engineer": 1,
			"manager": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1,
			"manager": 1,
			"senior": 1
		    }
		]
	    ],
	    "Periods": [     // This is also a four-year career path.
		1,
		1,
		1,
		1
	    ]
	}
    ]


The following command learns 5 hidden states from above training data
and prints the learned model in JSON format.


    wyi@u64b-> $GOPATH/bin/trainer -corpus=./toy_career.json -states=5 -iter=2000
    2014/10/07 12:53:25 Cannot create , output to stdout.
    {
      "S1": [
	0,
	0,
	2,   // Here the model says that both careers start with hidden state 2, 
	0,   // which, according to "Σγo", has weights on "senior", "staff", 
	0    // "software", "engineer".  This means that these words/roles are 
      ],     // clustered as they look similar (all contains "engineer").
      "S1Sum": 2,
      "Σγ": [
	1,
	1,
	2,
	2,
	0
      ],
      "Σξ": [
	[
	  0,
	  0,
	  1,
	  0,
	  0
	],
	[
	  0,
	  0,
	  0,
	  0,
	  1
	],
	[
	  0,
	  0,
	  0,
	  2,
	  0
	],
	[
	  1,
	  1,
	  0,
	  0,
	  0
	],
	[
	  0,
	  0,
	  0,
	  0,
	  0
	]
      ],
      "Σγo": [
	[
	  {
	    "Hist": {
	      "internet": 1
	    },
	    "Sum": 1
	  },
	  {
	    "Hist": {
	      "engineer": 1,          "software": 1,          "staff": 1
	    },
	    "Sum": 3
	  }
	],
	[
	  {
	    "Hist": {
	      "internet": 1
	    },
	    "Sum": 1
	  },
	  {
	    "Hist": {
	      "engineer": 1,          "manager": 1
	    },
	    "Sum": 2
	  }
	],
	[
	  {
	    "Hist": {
	      "internet": 3
	    },
	    "Sum": 3
	  },
	  {
	    "Hist": {
	      "engineer": 3,          "senior": 1,          "software": 3,          "staff": 1
	    },
	    "Sum": 8
	  }
	],
	[
	  {
	    "Hist": {
	      "internet": 2
	    },
	    "Sum": 2
	  },
	  {
	    "Hist": {
	      "engineer": 2,          "senior": 2,          "software": 2
	    },
	    "Sum": 6
	  }
	],
	[
	  {
	    "Hist": {
	      "internet": 1
	    },
	    "Sum": 1
	  },
	  {
	    "Hist": {
	      "engineer": 1,          "manager": 1,          "senior": 1
	    },
	    "Sum": 3
	  }
	]
      ]
    }


Conclusion, titles/positions/seniarity should be represented by a
certain level; instead by text (set of words).