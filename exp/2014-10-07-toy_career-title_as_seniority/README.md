
Following the [previous
experiment](../2014-10-07-toy_career-title_as_text/), which represent
seniority title by words, this experiment use seniority level (each
level is a unique text string).

Training corpus is as follows:

    [
	{
	    "Obs": [
		[
		    {
			"internet": 1
		    },
		    {
			"engineer": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"senior engineer": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"staff engineer": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"senior staff engineer": 1
		    }
		]
	    ],
	    "Periods": [
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
			"engineer": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"senior engineer": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"manager": 1
		    }
		],
		[
		    {
			"internet": 1
		    },
		    {
			"senior manager": 1
		    }
		]
	    ],
	    "Periods": [
		1,
		1,
		1,
		1
	    ]
	}
    ]
 
The model is as following:

    {
      "S1": [
	0,
	0,
	1,
	0,
	0,
	1
      ],
      "S1Sum": 2,
      "Σγ": [
	1,
	2,
	1,
	1,
	0,
	1
      ],
      "Σξ": [
	[
	  1,
	  0,
	  0,
	  0,
	  0,
	  0
	],
	[
	  1,
	  0,
	  0,
	  1,
	  0,
	  0
	],
	[
	  0,
	  1,
	  0,
	  0,
	  0,
	  0
	],
	[
	  0,
	  0,
	  0,
	  0,
	  1,
	  0
	],
	[
	  0,
	  0,
	  0,
	  0,
	  0,
	  0
	],
	[
	  0,
	  1,
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
	      "internet": 2
	    },
	    "Sum": 2
	  },
	  {
	    "Hist": {
	      "senior staff engineer": 1,
	      "staff engineer": 1
	    },
	    "Sum": 2
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
	      "senior engineer": 2
	    },
	    "Sum": 2
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
	      "engineer": 1
	    },
	    "Sum": 1
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
	      "manager": 1
	    },
	    "Sum": 1
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
	      "senior manager": 1
	    },
	    "Sum": 1
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
	      "engineer": 1
	    },
	    "Sum": 1
	  }
	]
      ]
    }


## Conclusion.

It is successful in this experiment to avoid the problem in the
[previous experiment](../2014-10-07-toy_career-title_as_text/) that
all levels of engineers are clustered together.

However, we observed that "staff engineer" and "senior staff engineer"
are clustered together, which is not what we really want.

This might be due to the randomness and that EM algorithm is easy to
be trapped in local optimal esitmation.

I am going to change the training code to introduce either

1. tempered EM, or
2. multiple estimate and model selection

to help EM more robust in reaching ``good'' optimal estimates.