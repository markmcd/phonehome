var app = angular.module('phonehome', []);
app.config(function($sceDelegateProvider) {
    $sceDelegateProvider.resourceUrlWhitelist([
        'self',
        'http://gist-it.appspot.com/**'
    ]);
});
app.controller('searchCtrl', function($scope, $http) {
    $scope.query = '';
    $scope.search = function() {
        //var results = search($scope.query);
        //$scope.results = results;
        $http({method: 'GET', url: "search?q=" + encodeURIComponent($scope.query)}).
            success(function(data, status, headers, config) {
                if (data != "null") {
                    $scope.results = data;
                    window.setTimeout(function(){
                        Prism.highlightAll();
                    }, 100);
                }
            });
    };
    $scope.gitHubUrl = function(result) {
        return 'http://gist-it.appspot.com/github/' + result.User + '/' +
            result.Repo + '/blob/master/' + result.Path + '?slice=' +
            (result.Line - 3) + ':' + (result.Line + 3)
    }
});
