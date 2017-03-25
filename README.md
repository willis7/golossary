# golossary
A chatbot for storing and retrieving words relating to a specific subject, text, or dialect, with explanations


## Architecture
TODO

```
         +----------------+        
         |   Message      |        
         |   Dispatcher   |        
         +----------------+        
            /          \           
           /            \          
          /              \         
         /                \        
        /                  \       
+-----------+         +-----------+
| Performer |         | Performer |
|           |         |           |
+-----------+         +-----------+
```


## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D


## Credits

Golossary uses a number of open source projects to work properly:

* [Gorilla] - Provides an implementation for the WebSocket protocol defined in RFC 6455.
* [Viper] - Complete configuration solution for go applications including 12 factor apps.

## License

MIT


[Gorilla]: http://www.gorillatoolkit.org/pkg/websocket#ReadJSON
[Viper]: https://github.com/spf13/viper
