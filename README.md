# Validation

Extended validation for [validator](http://gopkg.in/validator.v1) in
[Go](http://golang.org).

## Usage

You can use this as a replacement for [validator](http://gopkg.in/validator.v1):

    type User struct {
        Name    string   `validate:"required,min=5"`
        Email   string   `validate:"required,email"`
        Address Address  `validate:"required,nested"`
    }

    func (u User) Validate() error {
        return validation.Validate(u)
    }

    type Address struct {
        City string `validate:"required,min=5"
    }

    func (a Address) Validate() error {
        return validation.Validate(a)
    }

or just for adding some "missing" constraints:

    func init() {
        // Validates a value or its length when non-zero.
        validator.SetValidationFunc("min", validation.Minimum)
        // Validates a string to be an email address.
        validator.SetValidationFunc("email", validation.Email)
        // Validates a value to be not nil, not zero-value, not empty.
        validator.SetValidationFunc("required", validation.Required)
        // Find more validation funcs in the source code...
    }

Find more details by reading the
[godoc](http://godoc.org/github.com/satisfeet/go-validation).

## License

Copyright 2014 Bodo Kaiser <i@bodokaiser.io>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
