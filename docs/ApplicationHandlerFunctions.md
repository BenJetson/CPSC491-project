# Application Handler Functions
​
Notes on setup for the application handler functions
​
- Before starting on the handler functions, the database drivers must be
  implemented
- The handler functions to implement come from the database drivers
- First step for application handler functions is to check to see who is sending
  the request
  - When receiving who sends the request, check to see if they have proper
    authentication
- Information needs to be decoded into JSON format (use the decoder)
- Check the data that is sent over to ensure it is in a proper format
- Once authentication and formatting is checked, then either receive the data
  being sent over or send information to the database depending on the handler
  function being implemented
- Final step is to send the JSON response