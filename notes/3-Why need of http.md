- HTTP is an application layer protocol
- Underlaying protocol is TCP which is in tranport layer
- TCP just prodives "envelopes" that can transfer bytes around the network
- An application protocol assigns structure and meaning to the contents of the envelopes
- **TCP just makes sure the packet was received.** 

- HTTP can sent **Status Code (1xx to 5xx)** 
- You'd have to reimplement a whole status, notification, and state system yourself if you are just strictly working with raw packets.

- TCP doesn't know anything about the data that it passes around. 
- **Only 0s and 1s** received over TCP, We never know who sent a data, what data format is, what data is
- The only way around this is if there is an agreement/ shared convention between the applications running on the connected machines about what will happen
