## Disclaimer: 
The idea and the code provided here are for the sole purpose of security research. 
I am not responsiple for any malicious use of this repository or the api at `https://www.alkeldi.com/tracker`.

## Info
All the code is at https://github.com/alkeldi/outlook-tracker/blob/master/app/tracker.go,
the rest of this repository is just for using golang modules and docker.

# Idea
Microsoft's web version of the outlook mail-client provides a security check for emails. 
Whenever an email is opened, the client searches for all the links in the email and scans them for viruses.
At first, you would think that this behaviour is good for the user's security.
However, it can be easily exploited to violate the user's privacy. 

For the outlook client to be able to scan the links in an email, it has to visit those links. 
While the fetching is done on the backend by Microsoft's servers, 
it can still be used to track the count and the time whenever an email is opened or reopened by the receiver.

# Exploite
To write an exploite, all that's needed is a web server that's accessible by Microsoft servers (ie. connected to the internet). 
This webserver needs to generate a special link for tracking an email. 
This link must be included in the email that we want to track.

When the receiver opens our email, the outlook client will scan our link.
So, our  webserver will get a visit from a Microsoft server. We can store the information about this visit in a database. 
This database can be refrenced whenever we want to see how many times our email was opened or when it was opened.

# A working example
As a demo, I am currently running the code in this repository on https://www.alkeldi.com/tracker .
If you want to track an email you would do the following:
  1. First, create a new tracker by visiting https://www.alkeldi.com/tracker/create . This will show you something like `"YYNuzCr2FOLjzYX2Lyx4WPx7U5D4LZzVihrowpfDonNGgzQoXW"`
  This random string is the identifier of your tracker.
  2. In the email that you want to track, you need to include `https://www.alkeldi.com/tracker/track/YYNuzCr2FOLjzYX2Lyx4WPx7U5D4LZzVihrowpfDonNGgzQoXW`. 
  Replace the identifier with the one you got in step 1.
  3. Whenever the receiver opens your email, microsoft will visit the link in step 2, and the tracker will record this info.
  4. If you want to see the info about your tracker, visit https://www.alkeldi.com/tracker/visits/YYNuzCr2FOLjzYX2Lyx4WPx7U5D4LZzVihrowpfDonNGgzQoXW.
  Replace the identifier with the one you got in step 1.
  5. The first entry in the result of step 4 has an ip of -1.-1.-1.-1. This entry represents the time when the tracker was created.
  
# Important note
This was tested only on the webversion of the outlook client. the ios and android clients might have a different behaviour.
  
  









