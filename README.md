# Description

This is a golang repo that implemenents a (as of the time of writing) very inefficient, golang specific, sqlite function that invokes jev to classify arbitrary content. The function, `jev_classify`, implements (roughly) typesafe's classification API verbatim:

```
jev_classify(input, question, options)
```

Where `text` contains the thing-to-be-classified (ie, the `state`), `question` contains the classification instructions, and `options` contains the classification options. Here's an example invocation:

```
jev_classify("This cannot wait!", 'How urgent is this message?', json_object('urgent','needs attention right now','normal','routine, can wait','low','no action needed'))
```

And here's an example of the function in action on a real sqlite DB:
```
❯ sqlite3 tickets.db <<'EOF'
CREATE TABLE tickets (id INTEGER PRIMARY KEY, body TEXT NOT NULL);
INSERT INTO tickets (body) VALUES
  ('Production is down and customers are seeing 500 errors. Please help now!'),
  ('Our Stripe integration has been failing for 3 days and we are losing sales. ASAP please.'),
  ('Security alert: I think someone logged into my account from another country.'),
  ('Could you tell me when the next invoice will be sent? No rush.'),
  ('Is there a way to export my data as CSV? Just curious.'),
  ('How do I change the email address on my profile?'),
  ('Thanks for the quick fix last week, everything looks great now.'),
  ('Love the new dashboard design, nice work to the team!');
EOF

❯ go run main.go tickets.db "SELECT id, substr(body, 1, 40),     
  jev_classify(body, 'How urgent is this message?',
    json_object('urgent','needs attention right now','normal','routine, can wait','low','no action needed'))
  FROM tickets LIMIT 3"

1       Production is down and customers are see        urgent
2       Our Stripe integration has been failing         urgent
3       Security alert: I think someone logged i        urgent
```

The implementation is currently *very* inefficient - it makes a blocking call the jev API every time it's invoked, and doesn't use batching.

# Using it in your code
Nevertheless, if you want to use it, feel free; after importing the library, you can register a driver with the function available in your code like so:
```
sql.Register("sqlite3_jev", jevsqlite.Driver(os.Getenv("TYPESAFE_API_KEY")))
db, _ := sql.Open("sqlite3_jev", ...)
```
Happy classifying :)