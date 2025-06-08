package examples

func ExamplePostBodyCardBytes() []byte {
	return []byte(`
		{
			"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
			"type": "AdaptiveCard",
			"version": "1.0",
			"body": [
			{
				"type": "TextBlock",
				"text": "A disturbance in the force",
				"weight": "bolder",
				"size": "medium",
				"wrap": true
			},
			{
				"type": "ColumnSet",
				"columns": [
				{
					"type": "Column",
					"width": "auto",
					"items": [
					{
						"type": "Image",
						"url": "https://www4.pictures.zimbio.com/mp/ATCkWtsLsoEl.jpg",
						"size": "small",
						"style": "person"
					}
					]
				},
				{
					"type": "Column",
					"width": "stretch",
					"items": [
					{
						"type": "TextBlock",
						"text": "Obi-Wan Kenobi",
						"weight": "bolder",
						"wrap": true
					},
					{
						"type": "TextBlock",
						"spacing": "none",
						"text": "Created {{DATE(2017-02-14T06:08:39Z, SHORT)}}",
						"isSubtle": true,
						"wrap": true
					}
					]
				}
				]
			},
			{
				"type": "TextBlock",
				"text": "I felt something... as if millions of voices suddenly cried out in terror and were suddenly silenced.",
				"wrap": true
			},
			{
				"type": "FactSet",
				"facts": [
				{
					"title": "Current location:",
					"value": "On board the Millennium Falcon"
				},
				{
					"title": "Source:",
					"value": "Alderaan?"
				}
				]
			}
			]
		}
	`)
}
